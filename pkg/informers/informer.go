package informers

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	Start(stopCh <-chan struct{}) error
}

type InformerSet struct {
	Factories []dynamicinformer.DynamicSharedInformerFactory
	Resources []Resource
}

func (set *InformerSet) Start(stopCh <-chan struct{}) error {
	defer func() {
		for _, resource := range set.Resources {
			resource.Queue.ShutDown()
		}
	}()

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

	for _, resource := range set.Resources {
		for i := 0; i < resource.Workers; i++ {
			go wait.Until(resource.RunWorker, time.Second, stopCh)
		}
	}

	<-stopCh

	return nil
}

func Setup(options Options) (*InformerSet, error) {
	config := options.Config

	informerSet := &InformerSet{}

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
	klog.Infof("default channelNames: %v", defaultChannelNames)
	defaultWorkers := config.DefaultWorkers
	if defaultWorkers < 1 {
		defaultWorkers = 3
	}
	klog.Infof("default workers: %d", defaultWorkers)

	for i, namespace := range config.Namespaces {
		klog.Infof("[n%d] setup namespace %d/%d: %s",
			i+1, i+1, len(config.Namespaces), namespace.Namespace)

		namespaceDefaultResyncPeriodFunc, err := namespace.BuildResyncPeriodFuncWithDefault(
			defaultResyncPeriodFunc)
		if err != nil {
			return nil, errors.Wrapf(err,
				"namespace.BuildResyncPeriodFuncWithDefault error, .Namespaces[%d]: %s",
				i, namespace.MinResyncPeriod)
		}
		namespaceDefaultChannelNames := namespace.DefaultChannelNames
		if len(namespaceDefaultChannelNames) == 0 {
			namespaceDefaultChannelNames = defaultChannelNames
		}
		klog.Infof("[n%d] default channelNames: %v",
			i+1, namespaceDefaultChannelNames)
		namespaceDefaultWorkers := namespace.DefaultWorkers
		if namespaceDefaultWorkers < 1 {
			namespaceDefaultWorkers = defaultWorkers
		}
		klog.Infof("[n%d] default workers: %d",
			i+1, namespaceDefaultWorkers)

		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
			dynamicClient, namespaceDefaultResyncPeriodFunc(), namespace.Namespace, nil)

		for j, resource := range namespace.Resources {
			klog.Infof("[n%d,r%d] setup resource %d/%d: %s",
				i+1, j+1, j+1, len(namespace.Resources), resource.Resource)

			resyncPeriodFunc, err := resource.BuildResyncPeriodFuncWithDefault(
				namespaceDefaultResyncPeriodFunc)
			if err != nil {
				return nil, errors.Wrapf(err,
					"resource.BuildResyncPeriodFuncWithDefault error, .Namespaces[%d].Resources[%d]: %s",
					i, j, resource.ResyncPeriod)
			}
			channelNames := resource.ChannelNames
			if len(channelNames) == 0 {
				channelNames = namespaceDefaultChannelNames
			}
			klog.Infof("[n%d,r%d] final channelNames: %v", i+1, j+1, channelNames)
			workers := resource.Workers
			if workers < 1 {
				workers = namespaceDefaultWorkers
			}
			klog.Infof("[n%d,r%d] final workers: %v", i+1, j+1, workers)

			rateLimiter := workqueue.DefaultControllerRateLimiter()
			queue := workqueue.NewRateLimitingQueue(rateLimiter)
			gvr, _ := schema.ParseResourceArg(resource.Resource)
			if gvr == nil {
				return nil, errors.Wrapf(err,
					"schema.ParseResourceArg error, .Namespaces[%d].Resource[%d]: %s",
					i, j, resource.Resource)
			}
			informer := factory.ForResource(*gvr).Informer()
			handlerFuncs := cache.ResourceEventHandlerFuncs{}
			if resource.NoticeWhenAdded {
				klog.Infof("[n%d,r%d] set AddFunc", i+1, j+1)
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
			if resource.NoticeWhenDeleted {
				klog.Infof("[n%d,r%d] set DeleteFunc", i+1, j+1)
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
			if resource.NoticeWhenUpdated {
				klog.Infof("[%d,%d] set UpdateFunc", i+1, j+1)
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

			informerSet.Resources = append(informerSet.Resources, Resource{
				Resource:        resource.Resource,
				UpdateOn:        resource.UpdateOn,
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
