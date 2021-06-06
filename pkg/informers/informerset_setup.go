package informers

import (
	"github.com/pkg/errors"
	spg "github.com/spongeprojects/client-go/api/spongeprojects.com/v1alpha1"
	spgc "github.com/spongeprojects/client-go/client/clientset/versioned"
	spgi "github.com/spongeprojects/client-go/client/informers/externalversions"
	spgl "github.com/spongeprojects/client-go/client/listers/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/resourcebuilder"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"time"
)

const channelNamePrintToStdout = "print-to-stdout"

func Setup(config Config) (*InformerSet, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	defaultWorkers := config.DefaultWorkers
	if defaultWorkers < 1 {
		defaultWorkers = 3
	}
	defaultMaxRetries := config.DefaultMaxRetries
	if defaultMaxRetries < 1 {
		defaultMaxRetries = 3
	}
	defaultChannelNames := config.DefaultChannelNames
	defaultResyncPeriodFunc := buildResyncPeriodFuncByDuration(config.MinResyncPeriod)

	klog.V(1).Infof(
		"default: workers: %d, max retries: %d, channel names: %s",
		defaultWorkers, defaultMaxRetries, defaultChannelNames)

	restConfig, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	resourceBuilder, err := resourcebuilder.New(config.Kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "resourcebuilder.New error")
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create dynamic client error")
	}

	spgClientset, err := spgc.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "new clientset error")
	}

	spgInformerFactory := spgi.NewSharedInformerFactory(spgClientset, 12*time.Hour)

	var channelMap channels.ChannelMap
	var channelQueue workqueue.RateLimitingInterface
	var channelInformer cache.SharedIndexInformer
	var channelLister spgl.ChannelLister

	if config.JustWatch {
		printToStdout, _ := channels.NewChannelPrint(&spg.ChannelPrintConfig{
			Writer: channels.PrintWriterStdout,
		})
		channelMap = map[string]channels.Channel{
			channelNamePrintToStdout: printToStdout,
		}
	} else {
		channelMap = make(channels.ChannelMap)

		channelsRateLimiter := workqueue.DefaultControllerRateLimiter()
		channelQueue = workqueue.NewRateLimitingQueue(channelsRateLimiter)
		channelInformer = spgInformerFactory.Spongeprojects().V1alpha1().Channels().Informer()
		channelLister = spgInformerFactory.Spongeprojects().V1alpha1().Channels().Lister()

		channelInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				channel, ok := obj.(*spg.Channel)
				if !ok {
					return
				}
				klog.V(2).Infof("[channel] received: new channel of type %s: %s", channel.Spec.Type, channel.Name)
				channelQueue.Add(channel.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_, ok1 := oldObj.(*spg.Channel)
				channel, ok2 := newObj.(*spg.Channel)
				if !ok1 || !ok2 {
					return
				}
				klog.V(2).Infof("[channel] received: channel updated: %s", channel.Name)
				channelQueue.Add(channel.Name)
			},
			DeleteFunc: func(obj interface{}) {
				channel, ok := obj.(*spg.Channel)
				if !ok {
					return
				}
				klog.V(2).Infof("[channel] received: channel deleted: %s", channel.Name)
				channelQueue.Add(channel.Name)
			},
		})
	}

	watcherRateLimiter := workqueue.DefaultControllerRateLimiter()
	watcherQueue := workqueue.NewRateLimitingQueue(watcherRateLimiter)
	watcherInformer := spgInformerFactory.Spongeprojects().V1alpha1().Watchers().Informer()
	watcherLister := spgInformerFactory.Spongeprojects().V1alpha1().Watchers().Lister()

	watcherInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			watcher, ok := obj.(*spg.Watcher)
			if !ok {
				return
			}
			k, _ := cache.MetaNamespaceKeyFunc(watcher)
			klog.V(2).Infof("[watcher] received: new watcher for resource %s: %s", watcher.Spec.Resource, k)
			watcherQueue.Add(k)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// TODO: make it mutable
			//_, ok1 := oldObj.(*spg.Watcher)
			//watcher, ok2 := newObj.(*spg.Watcher)
			//if !ok1 || !ok2 {
			//	return
			//}
			//k, _ := cache.MetaNamespaceKeyFunc(watcher)
			//klog.V(5).Infof("[watcher] received: watcher updated: %s", k)
			//watcherQueue.Add(k)
		},
		DeleteFunc: func(obj interface{}) {
			watcher, ok := obj.(*spg.Watcher)
			if !ok {
				return
			}
			k, _ := cache.MetaNamespaceKeyFunc(watcher)
			klog.V(2).Infof("[watcher] received: watcher deleted: %s", k)
			watcherQueue.Add(k)
		},
	})

	clusterWatcherRateLimiter := workqueue.DefaultControllerRateLimiter()
	clusterWatcherQueue := workqueue.NewRateLimitingQueue(clusterWatcherRateLimiter)
	clusterWatcherInformer := spgInformerFactory.Spongeprojects().V1alpha1().ClusterWatchers().Informer()
	clusterWatcherLister := spgInformerFactory.Spongeprojects().V1alpha1().ClusterWatchers().Lister()

	clusterWatcherInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			watcher, ok := obj.(*spg.ClusterWatcher)
			if !ok {
				return
			}
			klog.V(2).Infof("[clusterwatcher] received: new cluster watcher for resource %s: %s",
				watcher.Spec.Resource, watcher.Name)
			clusterWatcherQueue.Add(watcher.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// TODO: make it mutable
			//_, ok1 := oldObj.(*spg.ClusterWatcher)
			//watcher, ok2 := newObj.(*spg.ClusterWatcher)
			//if !ok1 || !ok2 {
			//	return
			//}
			//klog.V(5).Infof("[clusterwatcher] received: cluster watcher updated: %s", watcher.Name)
			//clusterWatcherQueue.Add(watcher.Name)
		},
		DeleteFunc: func(obj interface{}) {
			watcher, ok := obj.(*spg.ClusterWatcher)
			if !ok {
				return
			}
			klog.V(2).Infof("[clusterwatcher] received: cluster watcher deleted: %s", watcher.Name)
			clusterWatcherQueue.Add(watcher.Name)
		},
	})

	return &InformerSet{
		JustWatch:               config.JustWatch,
		EventStore:              config.EventStore,
		DefaultWorkers:          defaultWorkers,
		DefaultMaxRetries:       defaultMaxRetries,
		DefaultChannelNames:     defaultChannelNames,
		DefaultResyncPeriodFunc: defaultResyncPeriodFunc,
		ChannelQueue:            channelQueue,
		ChannelInformer:         channelInformer,
		ChannelLister:           channelLister,
		ChannelMap:              channelMap,
		WatcherQueue:            watcherQueue,
		WatcherInformer:         watcherInformer,
		WatcherLister:           watcherLister,
		WatcherMap:              make(map[string]*Informer),
		ClusterWatcherQueue:     clusterWatcherQueue,
		ClusterWatcherInformer:  clusterWatcherInformer,
		ClusterWatcherLister:    clusterWatcherLister,
		ClusterWatcherMap:       make(map[string]*Informer),
		ResourceBuilder:         resourceBuilder,
		DynamicClient:           dynamicClient,
	}, nil
}
