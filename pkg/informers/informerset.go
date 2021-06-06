package informers

import (
	spgl "github.com/spongeprojects/client-go/client/listers/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/resourcebuilder"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"time"
)

type Interface interface {
	// Start starts all informers registered,
	// Start is non-blocking, you should always call Shutdown before exit.
	Start(stopCh <-chan struct{}) error

	// Shutdown should be called before exit.
	Shutdown()
}

type InformerSet struct {
	JustWatch  bool
	EventStore event_store.Interface

	DefaultWorkers          int
	DefaultMaxRetries       int
	DefaultChannelNames     []string
	DefaultResyncPeriodFunc ResyncPeriodFunc

	// ChannelQueue is the queue for channel delta, item: channel name
	ChannelQueue    workqueue.RateLimitingInterface
	ChannelInformer cache.SharedIndexInformer
	ChannelLister   spgl.ChannelLister
	ChannelMap      channels.ChannelMap

	// WatcherQueue is the queue for channel delta, item: watcher namespaced key
	WatcherQueue    workqueue.RateLimitingInterface
	WatcherInformer cache.SharedIndexInformer
	WatcherLister   spgl.WatcherLister

	// WatcherMap maps from namespaced key to *Informer
	WatcherMap map[string]*Informer

	// ClusterWatcherQueue is the queue for channel delta, item: cluster watcher name
	ClusterWatcherQueue    workqueue.RateLimitingInterface
	ClusterWatcherInformer cache.SharedIndexInformer
	ClusterWatcherLister   spgl.ClusterWatcherLister

	// ClusterWatcherMap maps from namespaced key to *Informer
	ClusterWatcherMap map[string]*Informer

	ResourceBuilder resourcebuilder.Interface
	DynamicClient   dynamic.Interface
}

func (s *InformerSet) Start(stopCh <-chan struct{}) error {
	if !s.JustWatch {
		go s.ChannelInformer.Run(stopCh)
	}
	go s.WatcherInformer.Run(stopCh)
	go s.ClusterWatcherInformer.Run(stopCh)

	klog.Info("waiting for caches to sync...")

	if !s.JustWatch {
		cache.WaitForCacheSync(stopCh, s.ChannelInformer.HasSynced)
	}
	cache.WaitForCacheSync(stopCh, s.WatcherInformer.HasSynced)
	cache.WaitForCacheSync(stopCh, s.ClusterWatcherInformer.HasSynced)

	klog.Info("caches synced, starting workers...")

	if !s.JustWatch {
		for i := 0; i < 3; i++ {
			go wait.Until(s.RunChannelWorker, time.Second, stopCh)
		}
	}

	for i := 0; i < 3; i++ {
		go wait.Until(s.RunWatcherWorker, time.Second, stopCh)
	}

	for i := 0; i < 3; i++ {
		go wait.Until(s.RunClusterWatcherWorker, time.Second, stopCh)
	}

	klog.Info("all workers started")

	<-stopCh

	return nil
}

func (s *InformerSet) Shutdown() {
	s.ChannelQueue.ShutDown()
	s.WatcherQueue.ShutDown()
	s.ClusterWatcherQueue.ShutDown()

	for _, informer := range s.WatcherMap {
		informer.ShutDown()
	}

	for _, informer := range s.ClusterWatcherMap {
		informer.ShutDown()
	}
}
