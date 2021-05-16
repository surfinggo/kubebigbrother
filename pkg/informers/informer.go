package informers

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

type Options struct {
	KubeConfig string

	Config *Config
}

type Interface interface {
	// Start starts all informers registered,
	// Start is non-blocking, you should always call Shutdown before exit.
	Start(stopCh <-chan struct{}) error

	// Shutdown should be called before exit.
	Shutdown()
}

type InformerSet struct {
	ResourceBuilder   *ResourceBuilder
	Factories         []dynamicinformer.DynamicSharedInformerFactory
	ResourceInformers []Resource
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

	for _, resourceInformer := range set.ResourceInformers {
		for i := 0; i < resourceInformer.Workers; i++ {
			go wait.Until(resourceInformer.RunWorker, time.Second, stopCh)
		}
	}

	<-stopCh

	return nil
}

func (set *InformerSet) Shutdown() {
	for _, resourceInformer := range set.ResourceInformers {
		resourceInformer.Queue.ShutDown()
	}
}

func Setup(options Options) (*InformerSet, error) {
	config := options.Config

	informerSet := &InformerSet{}

	resourceBuilder, err := NewResourceBuilder(options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "NewResourceBuilder error")
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

	channelMap := make(ChannelMap)
	for name, channelConfig := range config.Channels {
		channel, err := BuildChannelFromConfig(&channelConfig)
		if err != nil {
			return nil, errors.Wrap(err, "build channel error")
		}
		channelMap[name] = channel
	}

	defaultResyncPeriodFunc, err := config.BuildResyncPeriodFunc()
	if err != nil {
		return nil, errors.Wrap(err, "config.BuildResyncPeriodFunc error")
	}
	defaultChannelNames := config.DefaultChannelNames
	klog.V(1).Infof("default channelNames: %v", defaultChannelNames)
	defaultWorkers := config.DefaultWorkers
	if defaultWorkers < 1 {
		defaultWorkers = 3
	}
	klog.V(1).Infof("default workers: %d", defaultWorkers)

	for i, namespaceConfig := range config.Namespaces {
		klog.Infof("[n%d] setup namespace %d/%d: %s",
			i, i+1, len(config.Namespaces), namespaceConfig.Namespace)

		namespaceDefaultResyncPeriodFunc, err := namespaceConfig.BuildResyncPeriodFuncWithDefault(
			defaultResyncPeriodFunc)
		if err != nil {
			return nil, errors.Wrapf(err,
				"namespace.BuildResyncPeriodFuncWithDefault error, .Namespaces[%d]: %s",
				i, namespaceConfig.MinResyncPeriod)
		}
		namespaceDefaultChannelNames := namespaceConfig.DefaultChannelNames
		if len(namespaceDefaultChannelNames) == 0 {
			namespaceDefaultChannelNames = defaultChannelNames
		}
		klog.V(1).Infof("[n%d] default channelNames: %v",
			i, namespaceDefaultChannelNames)
		namespaceDefaultWorkers := namespaceConfig.DefaultWorkers
		if namespaceDefaultWorkers < 1 {
			namespaceDefaultWorkers = defaultWorkers
		}
		klog.V(1).Infof("[n%d] default workers: %d",
			i, namespaceDefaultWorkers)

		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
			dynamicClient, namespaceDefaultResyncPeriodFunc(), namespaceConfig.Namespace, nil)

		duplicate := make(map[string]bool)

		for j, resourceConfig := range namespaceConfig.Resources {
			klog.Infof("[n%d,r%d] setup resource %d/%d: %s",
				i, j, j+1, len(namespaceConfig.Resources), resourceConfig.Resource)

			if _, ok := duplicate[resourceConfig.Resource]; ok {
				return nil, errors.Errorf(
					"duplicated resources in same namespace, .Namespaces[%d].Resources[%d]: %s",
					i, j, resourceConfig.Resource)
			}
			duplicate[resourceConfig.Resource] = true

			resyncPeriodFunc, err := resourceConfig.BuildResyncPeriodFuncWithDefault(
				namespaceDefaultResyncPeriodFunc)
			if err != nil {
				return nil, errors.Wrapf(err,
					"resource.BuildResyncPeriodFuncWithDefault error, .Namespaces[%d].Resources[%d]: %s",
					i, j, resourceConfig.ResyncPeriod)
			}
			channelNames := resourceConfig.ChannelNames
			if len(channelNames) == 0 {
				channelNames = namespaceDefaultChannelNames
			}
			klog.V(1).Infof("[n%d,r%d] final channelNames: %v", i, j, channelNames)
			workers := resourceConfig.Workers
			if workers < 1 {
				workers = namespaceDefaultWorkers
			}
			klog.V(1).Infof("[n%d,r%d] final workers: %v", i, j, workers)

			rateLimiter := workqueue.DefaultControllerRateLimiter()
			queue := workqueue.NewRateLimitingQueue(rateLimiter)

			gvr, err := informerSet.ResourceBuilder.ParseGroupResource(resourceConfig.Resource)
			if err != nil {
				return nil, errors.Wrapf(err,
					"parse resource error, .Namespaces[%d].Resource[%d]: %s",
					i, j, resourceConfig.Resource)
			}

			informer := factory.ForResource(gvr).Informer()
			handlerFuncs := cache.ResourceEventHandlerFuncs{}
			if resourceConfig.NoticeWhenAdded {
				klog.V(1).Infof("[n%d,r%d] set AddFunc", i, j)
				handlerFuncs.AddFunc = func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					e := &Event{
						Type:         EventTypeAdded,
						Obj:          s,
						ChannelNames: channelNames,
					}
					klog.V(5).Infof("received: [%s] %s", e.Type, NamespaceKey(s))
					queue.Add(e)
				}
			}
			if resourceConfig.NoticeWhenDeleted {
				klog.V(1).Infof("[n%d,r%d] set DeleteFunc", i, j)
				handlerFuncs.DeleteFunc = func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					e := &Event{
						Type:         EventTypeDeleted,
						Obj:          s,
						ChannelNames: channelNames,
					}
					klog.V(5).Infof("received: [%s] %s", e.Type, NamespaceKey(s))
					queue.Add(e)
				}
			}
			if resourceConfig.NoticeWhenUpdated {
				klog.V(1).Infof("[%d,%d] set UpdateFunc", i, j)
				handlerFuncs.UpdateFunc = func(oldObj, obj interface{}) {
					oldS, ok1 := oldObj.(*unstructured.Unstructured)
					s, ok2 := obj.(*unstructured.Unstructured)
					if !ok1 || !ok2 {
						return
					}
					e := &Event{
						Type:         EventTypeUpdated,
						Obj:          s,
						OldObj:       oldS,
						ChannelNames: channelNames,
					}
					klog.V(5).Infof("received: [%s] %s", e.Type, NamespaceKey(s))
					queue.Add(e)
				}
			}
			informer.AddEventHandlerWithResyncPeriod(handlerFuncs, resyncPeriodFunc())

			informerSet.ResourceInformers = append(informerSet.ResourceInformers, Resource{
				Resource:        resourceConfig.Resource,
				UpdateOn:        resourceConfig.UpdateOn,
				ChannelMap:      channelMap,
				Queue:           queue,
				Workers:         workers,
				processingItems: &sync.WaitGroup{}, // TODO: wait before exit
			})
		} // end resources loop

		// a factory for a namespace
		informerSet.Factories = append(informerSet.Factories, factory)
	} // end namespaces loop

	return informerSet, nil
}
