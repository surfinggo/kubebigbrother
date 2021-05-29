package informers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/resourcebuilder"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Options struct {
	KubeConfig string

	Config *Config
}

func (o *Options) Validate() error {
	if err := o.Config.Validate(); err != nil {
		return errors.Wrap(err, "invalid config")
	}
	return nil
}

type Interface interface {
	// Start starts all informers registered,
	// Start is non-blocking, you should always call Shutdown before exit.
	Start(stopCh <-chan struct{}) error

	// Shutdown should be called before exit.
	Shutdown()
}

type InformerSet struct {
	ResourceBuilder   resourcebuilder.Interface
	Factories         []dynamicinformer.DynamicSharedInformerFactory
	ResourceInformers []*Informer
}

func (set *InformerSet) Start(stopCh <-chan struct{}) error {
	for i, factory := range set.Factories {
		klog.Infof("starting factory %d/%d", i+1, len(set.Factories))
		go factory.Start(stopCh)
	}

	for i, factory := range set.Factories {
		for gvr, ok := range factory.WaitForCacheSync(stopCh) {
			if !ok {
				return fmt.Errorf(
					"timed out waiting for caches to sync, .Factories[%d][%s]", i, gvr)
			}
		}
	}

	for i, resourceInformer := range set.ResourceInformers {
		klog.Infof("starting informer %d/%d: %s, workers: %d, resource: %s",
			i+1, len(set.ResourceInformers), resourceInformer.ID,
			resourceInformer.Workers, resourceInformer.Resource)
		for i := 0; i < resourceInformer.Workers; i++ {
			go wait.Until(resourceInformer.RunWorker, time.Second, stopCh)
		}
	}

	<-stopCh

	return nil
}

func (set *InformerSet) Shutdown() {
	for _, resourceInformer := range set.ResourceInformers {
		resourceInformer.ShutDown()
	}
}

// Setup setups new InformerSet
func Setup(options Options) (*InformerSet, error) {
	if err := options.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid options")
	}

	config := options.Config

	informerSet := &InformerSet{}

	resourceBuilder, err := resourcebuilder.New(options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "resourcebuilder.New error")
	}
	informerSet.ResourceBuilder = resourceBuilder

	restConfig, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "dynamic.NewForConfig error")
	}

	channelMap := make(channels.ChannelMap)
	for name, channelConfig := range config.Channels {
		channel, err := setupChannelFromConfig(&channelConfig)
		if err != nil {
			return nil, errors.Wrap(err, "setup channel error")
		}
		channelMap[name] = channel
	}

	// build []ChannelToProcess for an event
	buildChannelsToProcess := func(e *event.Event, names []channels.ChannelName) []ChannelToProcess {
		var channelsToProcess []ChannelToProcess
		for _, name := range names {
			if channel, ok := channelMap[name]; ok {
				channelsToProcess = append(channelsToProcess, ChannelToProcess{
					ChannelName:         name,
					EventProcessContext: channel.NewEventProcessContext(e),
				})
			}
		}
		return channelsToProcess
	}

	defaultResyncPeriodFunc, err := config.buildResyncPeriodFunc()
	if err != nil {
		return nil, errors.Wrap(err, "config.BuildResyncPeriodFunc error")
	}
	defaultChannelNames := config.DefaultChannelNames
	defaultWorkers := config.DefaultWorkers
	if defaultWorkers < 1 {
		defaultWorkers = 3
	}
	defaultMaxRetries := config.DefaultMaxRetries
	if defaultMaxRetries < 1 {
		defaultMaxRetries = 3
	}
	klog.V(1).Infof(
		"global default: workers: %d, max retries: %d, channel names: %s",
		defaultWorkers, defaultMaxRetries, defaultChannelNames)

	for i, namespaceConfig := range config.Namespaces {
		nID := fmt.Sprintf("n%d", i) // unique id
		nDesc := fmt.Sprintf(".Namespaces[%d]", i)

		klog.Infof("[%s] setup namespace %d/%d: %s",
			nID, i+1, len(config.Namespaces), namespaceConfig.Namespace)

		namespaceDefaultResyncPeriodFunc, err := namespaceConfig.buildResyncPeriodFuncWithDefault(
			defaultResyncPeriodFunc)
		if err != nil {
			return nil, errors.Wrapf(err,
				"invalid resync period in %s: %s", nDesc, namespaceConfig.MinResyncPeriod)
		}
		namespaceDefaultChannelNames := namespaceConfig.DefaultChannelNames
		if len(namespaceDefaultChannelNames) == 0 {
			namespaceDefaultChannelNames = defaultChannelNames
		}
		namespaceDefaultWorkers := namespaceConfig.DefaultWorkers
		if namespaceDefaultWorkers < 1 {
			namespaceDefaultWorkers = defaultWorkers
		}
		namespaceDefaultMaxRetries := namespaceConfig.DefaultMaxRetries
		if namespaceDefaultMaxRetries < 1 {
			namespaceDefaultMaxRetries = defaultMaxRetries
		}
		klog.V(1).Infof(
			"[%s] namespace default: workers: %d, max retries: %d, channel names: %v",
			nID, namespaceDefaultWorkers,
			namespaceDefaultMaxRetries,
			namespaceDefaultChannelNames)

		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
			dynamicClient, namespaceDefaultResyncPeriodFunc(),
			namespaceConfig.Namespace, nil)

		duplicate := make(map[string]bool)

		for j, resourceConfig := range namespaceConfig.Resources {
			rID := fmt.Sprintf("%sr%d", nID, j) // unique rID
			rDesc := fmt.Sprintf("%s.Resources[%d]", nDesc, i)

			klog.Infof("[%s] setup resource %d/%d: %s",
				rID, j+1, len(namespaceConfig.Resources), resourceConfig.Resource)

			gvr, err := informerSet.ResourceBuilder.ParseGroupResource(resourceConfig.Resource)
			if err != nil {
				return nil, errors.Wrapf(err,
					"invalid resource in %s: %s", rDesc, resourceConfig.Resource)
			}

			if _, ok := duplicate[gvr.String()]; ok {
				return nil, errors.Errorf(
					"duplicated resources in same namespace %s: %s",
					rDesc, resourceConfig.Resource)
			}
			duplicate[gvr.String()] = true

			resyncPeriodFunc, err := resourceConfig.buildResyncPeriodFuncWithDefault(
				namespaceDefaultResyncPeriodFunc)
			if err != nil {
				return nil, errors.Wrapf(err,
					"invalid resync period in %s: %s",
					rDesc, resourceConfig.ResyncPeriod)
			}
			channelNames := resourceConfig.ChannelNames
			if len(channelNames) == 0 {
				channelNames = namespaceDefaultChannelNames
			}
			workers := resourceConfig.Workers
			if workers < 1 {
				workers = namespaceDefaultWorkers
			}
			maxRetries := resourceConfig.MaxRetries
			if maxRetries < 1 {
				maxRetries = namespaceDefaultMaxRetries
			}
			klog.V(1).Infof(
				"[%s] gvr: [%v], workers: %d, max retries: %d, channel names: %v",
				rID, gvr, workers, maxRetries, channelNames)

			rateLimiter := workqueue.DefaultControllerRateLimiter()
			queue := workqueue.NewRateLimitingQueue(rateLimiter)

			informer := factory.ForResource(gvr).Informer()
			handlerFuncs := cache.ResourceEventHandlerFuncs{}
			if !resourceConfig.NoticeWhenAdded &&
				!resourceConfig.NoticeWhenDeleted &&
				!resourceConfig.NoticeWhenUpdated {
				// for some reason, informer won't start if they are all false,
				// maybe someday someone can clarify why this happen.
				return nil, errors.Errorf(
					"NoticeWhenAdded, NoticeWhenDeleted and NoticeWhenUpdated "+
						"cannot be false simultaneously in %s", rDesc)
			}
			if resourceConfig.NoticeWhenAdded {
				klog.V(1).Infof("[%s] set AddFunc", rID)
				handlerFuncs.AddFunc = func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					e := event.NewAdded(s)
					klog.V(5).Infof("received: [%s] [%s]", e.Type, utils.GroupVersionKindName(s))
					queue.Add(&eventWrapper{
						Event:             e,
						ChannelsToProcess: buildChannelsToProcess(e, channelNames),
					})
				}
			}
			if resourceConfig.NoticeWhenDeleted {
				klog.V(1).Infof("[%s] set DeleteFunc", rID)
				handlerFuncs.DeleteFunc = func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					e := event.NewDeleted(s)
					klog.V(5).Infof("received: [%s] [%s]", e.Type, utils.GroupVersionKindName(s))
					queue.Add(&eventWrapper{
						Event:             e,
						ChannelsToProcess: buildChannelsToProcess(e, channelNames),
					})
				}
			}
			if resourceConfig.NoticeWhenUpdated {
				klog.V(1).Infof("[%s] set UpdateFunc", rID)
				handlerFuncs.UpdateFunc = func(oldObj, obj interface{}) {
					oldS, ok1 := oldObj.(*unstructured.Unstructured)
					s, ok2 := obj.(*unstructured.Unstructured)
					if !ok1 || !ok2 {
						return
					}
					updated := false
					for _, field := range resourceConfig.UpdateOn {
						f1, exist1, err1 := unstructured.NestedFieldNoCopy(
							s.Object, strings.Split(strings.TrimPrefix(field, "."), ".")...)
						f2, exist2, err2 := unstructured.NestedFieldNoCopy(
							oldS.Object, strings.Split(strings.TrimPrefix(field, "."), ".")...)
						if !exist1 || !exist2 {
							klog.Warningf("field not exist in resource: %s: %s",
								utils.GroupVersionKindName(s), field)
						}
						if err1 != nil {
							klog.Warningf("get field value error, resource: %s: %s",
								utils.GroupVersionKindName(s), err1)
						}
						if err2 != nil {
							klog.Warningf("get field value error, resource: %s: %s",
								utils.GroupVersionKindName(s), err2)
						}
						if exist1 && exist2 && err1 == nil && err2 == nil && !reflect.DeepEqual(f1, f2) {
							updated = true
							break
						}
					}
					if resourceConfig.UpdateOn == nil || updated {
						e := event.NewUpdated(s, oldS)
						klog.V(5).Infof("received: [%s] [%s]", e.Type, utils.GroupVersionKindName(s))
						queue.Add(&eventWrapper{
							Event:             e,
							ChannelsToProcess: buildChannelsToProcess(e, channelNames),
						})
					}
				}
			}
			informer.AddEventHandlerWithResyncPeriod(handlerFuncs, resyncPeriodFunc())

			informerSet.ResourceInformers = append(informerSet.ResourceInformers, &Informer{
				ID:              rID,
				Resource:        resourceConfig.Resource,
				UpdateOn:        resourceConfig.UpdateOn,
				ChannelMap:      channelMap,
				Queue:           queue,
				Workers:         workers,
				MaxRetries:      maxRetries,
				processingItems: &sync.WaitGroup{}, // TODO: wait before exit
			})
		} // end resources loop

		// a factory for a namespace
		informerSet.Factories = append(informerSet.Factories, factory)
	} // end namespaces loop

	return informerSet, nil
}
