package informers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/resourcebuilder"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"reflect"
	"strings"
	"sync"
)

// Bootstrapper builds an InformerSet from scratch
type Bootstrapper struct {
	config *Config

	globalDefaultResyncPeriodFunc ResyncPeriodFunc
	globalDefaultChannelNames     []ChannelName
	globalDefaultWorkers          int
	globalDefaultMaxRetries       int

	// channelMap is a map from channel name to channel instance,
	channelMap channels.ChannelMap

	// buildChannelsToProcess is used to build []ChannelToProcess to an event
	buildChannelsToProcess func(e *event.Event, names []ChannelName) []ChannelToProcess

	resourceBuilder *resourcebuilder.ResourceBuilder

	dynamicClient dynamic.Interface

	duplicatedName map[string]bool
}

func (b *Bootstrapper) buildInformerSet() (*InformerSet, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	if err := b.initGlobalDefaults(); err != nil {
		return nil, err
	}

	if err := b.initChannelMap(); err != nil {
		return nil, err
	}

	if err := b.initResourceBuilder(); err != nil {
		return nil, err
	}

	if err := b.initDynamicClient(); err != nil {
		return nil, err
	}

	informerSet := &InformerSet{}

	informerSet.ResourceBuilder = b.resourceBuilder

	for i, namespaceConfig := range b.config.ConfigFile.Namespaces {
		namespaceID := fmt.Sprintf("n%d", i) // unique namespaceID
		namespaceDesc := fmt.Sprintf(".Namespaces[%d]", i)

		klog.V(1).Infof("[%s] setup namespace %d/%d: %s",
			namespaceID, i+1,
			len(b.config.ConfigFile.Namespaces),
			namespaceConfig.Namespace)

		factory, informers, err := b.buildNamespacedInformerFactory(
			namespaceID, namespaceDesc, namespaceConfig)
		if err != nil {
			return nil, err
		}

		informerSet.Factories = append(informerSet.Factories, factory)
		for _, informer := range informers {
			informerSet.ResourceInformers = append(
				informerSet.ResourceInformers, informer)
		}
	} // end namespaces loop

	return informerSet, nil
}

func (b *Bootstrapper) buildNamespacedInformerFactory(
	namespaceID, namespaceDesc string, c NamespaceConfig) (
	dynamicinformer.DynamicSharedInformerFactory, []*Informer, error) {

	namespaceDefaultResyncPeriodFunc, err := c.buildResyncPeriodFunc(
		b.globalDefaultResyncPeriodFunc)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"invalid resync period in %s: %s", namespaceDesc, c.MinResyncPeriod)
	}

	namespaceDefaultChannelNames := c.getDefaultChannelNames(
		b.globalDefaultChannelNames)
	namespaceDefaultWorkers := c.getDefaultWorkers(
		b.globalDefaultWorkers)
	namespaceDefaultMaxRetries := c.getDefaultMaxRetries(
		b.globalDefaultMaxRetries)

	klog.V(2).Infof(
		"[%s] namespace default: workers: %d, max retries: %d, channel names: %v",
		namespaceID, namespaceDefaultWorkers,
		namespaceDefaultMaxRetries,
		namespaceDefaultChannelNames)

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		b.dynamicClient, namespaceDefaultResyncPeriodFunc(),
		c.Namespace, nil)

	duplicated := make(map[string]bool)

	var informers []*Informer
	for i, resourceConfig := range c.Resources {
		resourceID := fmt.Sprintf("%sr%d", namespaceID, i) // unique resourceID
		resourceDesc := fmt.Sprintf("%s.Resources[%d]", namespaceDesc, i)

		klog.V(1).Infof("[%s] setup resource %d/%d: %s",
			resourceID, i+1, len(c.Resources), resourceConfig.Resource)

		informer, err := b.buildResourceInformer(
			factory,
			namespaceDefaultResyncPeriodFunc,
			namespaceDefaultChannelNames,
			namespaceDefaultWorkers,
			namespaceDefaultMaxRetries,
			resourceID, resourceDesc, resourceConfig)
		if err != nil {
			return nil, nil, err
		}

		gvr := informer.GVR.String()
		if _, ok := duplicated[gvr]; ok {
			return nil, nil, errors.Errorf(
				"duplicated resources in same namespace %s: %s (gvr: [%s])",
				resourceDesc, resourceConfig.Resource, gvr)
		}
		duplicated[gvr] = true

		informers = append(informers, informer)
	} // end resources loop

	return factory, informers, nil
}

func (b *Bootstrapper) buildResourceInformer(
	factory dynamicinformer.DynamicSharedInformerFactory,
	namespaceDefaultResyncPeriodFunc ResyncPeriodFunc,
	namespaceDefaultChannelNames []ChannelName,
	namespaceDefaultWorkers int,
	namespaceDefaultMaxRetries int,
	resourceID string,
	resourceDesc string,
	c ResourceConfig) (*Informer, error) {
	if _, ok := b.duplicatedName[c.Name]; ok {
		return nil, errors.Errorf(
			"duplicated informer name %s: %s)", resourceDesc, c.Name)
	}
	b.duplicatedName[c.Name] = true

	gvr, err := b.resourceBuilder.ParseGroupResource(c.Resource)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid resource in %s: %s", resourceDesc, c.Resource)
	}

	saveSilently := func(e *event.Event) {
		b.config.EventStore.SaveSilently(e.ToModel(c.Name, gvr))
	}

	resyncPeriodFunc, err := c.buildResyncPeriodFunc(
		namespaceDefaultResyncPeriodFunc)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid resync period in %s: %s",
			resourceDesc, c.ResyncPeriod)
	}
	channelNames := c.getChannelNames(namespaceDefaultChannelNames)
	workers := c.getWorkers(namespaceDefaultWorkers)
	maxRetries := c.getMaxRetries(namespaceDefaultMaxRetries)
	klog.V(2).Infof(
		"[%s] gvr: [%v], workers: %d, max retries: %d, channel names: %v",
		resourceID, gvr, workers, maxRetries, channelNames)

	rateLimiter := workqueue.DefaultControllerRateLimiter()
	queue := workqueue.NewRateLimitingQueue(rateLimiter)

	informer := factory.ForResource(gvr).Informer()
	handlerFuncs := cache.ResourceEventHandlerFuncs{}
	if !c.NoticeWhenAdded &&
		!c.NoticeWhenDeleted &&
		!c.NoticeWhenUpdated {
		// for some reason, informer won't start if they are all false,
		// maybe someday someone can clarify why this happen.
		return nil, errors.Errorf(
			"NoticeWhenAdded, NoticeWhenDeleted and NoticeWhenUpdated "+
				"cannot be false simultaneously in %s", resourceDesc)
	}
	if c.NoticeWhenAdded {
		handlerFuncs.AddFunc = func(obj interface{}) {
			s, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			e := event.NewAdded(s)

			if b.config.SaveEvent {
				isCurrentlyAdded, err := b.config.EventStore.IsCurrentlyAdded(
					c.Name, gvr.Group, gvr.Version, gvr.Resource,
					e.Obj.GetNamespace(), e.Obj.GetName())
				if err != nil {
					klog.Warning(errors.Wrap(err, "find latest record error"))
				} else if isCurrentlyAdded {
					klog.V(5).Infof(
						"[%s] resource is currently added, skip ADDED event: [%s] [%s]",
						resourceID, e.Type, e.GroupVersionKindName())
					return // ADDED event are emitted when controller restart
				}
			}

			klog.V(5).Infof("[%s] received: [%s] [%s]",
				resourceID, e.Type, e.GroupVersionKindName())

			if b.config.SaveEvent {
				go saveSilently(e)
			}

			queue.Add(&eventWrapper{
				Event:             e,
				ChannelsToProcess: b.buildChannelsToProcess(e, channelNames),
			})
		}
	}
	if c.NoticeWhenDeleted {
		handlerFuncs.DeleteFunc = func(obj interface{}) {
			s, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			e := event.NewDeleted(s)

			klog.V(5).Infof("[%s] received: [%s] [%s]",
				resourceID, e.Type, utils.GroupVersionKindName(s))

			if b.config.SaveEvent {
				go saveSilently(e)
			}

			queue.Add(&eventWrapper{
				Event:             e,
				ChannelsToProcess: b.buildChannelsToProcess(e, channelNames),
			})
		}
	}
	if c.NoticeWhenUpdated {
		handlerFuncs.UpdateFunc = func(oldObj, obj interface{}) {
			oldS, ok1 := oldObj.(*unstructured.Unstructured)
			s, ok2 := obj.(*unstructured.Unstructured)
			if !ok1 || !ok2 {
				return
			}
			updated := false
			for _, field := range c.UpdateOn {
				fieldPath := strings.Split(strings.TrimPrefix(field, "."), ".")
				f1, exist1, err1 := unstructured.NestedFieldNoCopy(
					s.Object, fieldPath...)
				f2, exist2, err2 := unstructured.NestedFieldNoCopy(
					oldS.Object, fieldPath...)
				if !exist1 || !exist2 {
					klog.Warningf("[%s] field not exist in resource: %s: %s",
						resourceID, utils.GroupVersionKindName(s), field)
				}
				if err1 != nil {
					klog.Warningf("[%s] get field value error, resource: %s: %s",
						resourceID, utils.GroupVersionKindName(s), err1)
				}
				if err2 != nil {
					klog.Warningf("[%s] get field value error, resource: %s: %s",
						resourceID, utils.GroupVersionKindName(s), err2)
				}
				if exist1 && exist2 && err1 == nil && err2 == nil &&
					!reflect.DeepEqual(f1, f2) {
					updated = true
					break
				}
			}
			if c.UpdateOn == nil || updated {
				e := event.NewUpdated(s, oldS)

				klog.V(5).Infof("[%s] received: [%s] [%s]",
					resourceID, e.Type, utils.GroupVersionKindName(s))

				if b.config.SaveEvent {
					go saveSilently(e)
				}

				queue.Add(&eventWrapper{
					Event:             e,
					ChannelsToProcess: b.buildChannelsToProcess(e, channelNames),
				})
			}
		}
	}
	informer.AddEventHandlerWithResyncPeriod(handlerFuncs, resyncPeriodFunc())

	return &Informer{
		ID:              resourceID,
		Resource:        c.Resource,
		GVR:             gvr,
		UpdateOn:        c.UpdateOn,
		ChannelMap:      b.channelMap,
		Queue:           queue,
		Workers:         workers,
		MaxRetries:      maxRetries,
		processingItems: &sync.WaitGroup{}, // TODO: wait before exit
	}, nil
}

func (b *Bootstrapper) validate() error {
	return b.config.Validate()
}

func (b *Bootstrapper) initGlobalDefaults() error {
	resyncPeriodFunc, err := b.config.ConfigFile.buildResyncPeriodFunc()
	if err != nil {
		return errors.Wrap(err, "config.BuildResyncPeriodFunc error")
	}

	b.globalDefaultResyncPeriodFunc = resyncPeriodFunc
	b.globalDefaultChannelNames = b.config.ConfigFile.DefaultChannelNames
	b.globalDefaultWorkers = b.config.ConfigFile.getDefaultWorkers()
	b.globalDefaultMaxRetries = b.config.ConfigFile.getDefaultMaxRetries()

	klog.V(2).Infof(
		"global default: workers: %d, max retries: %d, channel names: %s",
		b.globalDefaultWorkers,
		b.globalDefaultMaxRetries,
		b.globalDefaultChannelNames)

	return nil
}

func (b *Bootstrapper) initChannelMap() error {
	i := 0
	b.channelMap = make(channels.ChannelMap)
	for name, channelConfig := range b.config.ConfigFile.Channels {
		i += 1
		klog.V(1).Infof("setup channel %d/%d: %s, type: %s",
			i, len(b.config.ConfigFile.Channels), name, channelConfig.Type)
		channel, err := channelConfig.setupChannel()
		if err != nil {
			return errors.Wrap(err, "setup channel error")
		}
		b.channelMap[name] = channel
	}

	b.buildChannelsToProcess = func(
		e *event.Event, names []ChannelName) (channelsToProcess []ChannelToProcess) {
		for _, name := range names {
			if channel, ok := b.channelMap[name]; ok {
				channelsToProcess = append(channelsToProcess, ChannelToProcess{
					ChannelName:         name,
					EventProcessContext: channel.NewEventProcessContext(e),
				})
			}
		}
		return channelsToProcess
	}

	return nil
}

func (b *Bootstrapper) initResourceBuilder() error {
	resourceBuilder, err := resourcebuilder.New(b.config.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "resourcebuilder.New error")
	}
	b.resourceBuilder = resourceBuilder
	return nil
}

func (b *Bootstrapper) initDynamicClient() error {
	restConfig, err := clientcmd.BuildConfigFromFlags("", b.config.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "get kube config error")
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return errors.Wrap(err, "create dynamic client error")
	}

	b.dynamicClient = dynamicClient
	return nil
}

// NewBootstrapper creates new Bootstrapper
func NewBootstrapper(config *Config) *Bootstrapper {
	return &Bootstrapper{
		config:         config,
		duplicatedName: make(map[string]bool),
	}
}
