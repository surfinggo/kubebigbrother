package informers

import (
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"time"
)

func (s *InformerSet) RunWatcherWorker() {
	for s.processNextWatcher() {
	}
}

// processNextWatcher waits and processes items in the queue
func (s *InformerSet) processNextWatcher() bool {
	// block until an key arrives or queue shutdown
	obj, shutdown := s.WatcherQueue.Get()
	if shutdown {
		return false
	}
	key := obj.(string)

	if klog.V(5).Enabled() {
		klog.Infof("[watcher] [%s try] key pop from queue: [%s]",
			humanize.Ordinal(s.WatcherQueue.NumRequeues(key)+1), key)
	}

	// we need to mark key as completed whether success or fail
	defer s.WatcherQueue.Done(key)

	result := s.processWatcher(key)
	s.handleWatcherErr(key, result)

	return true
}

// processWatcher process an item synchronously
func (s *InformerSet) processWatcher(key string) error {
	namespace, name, _ := cache.SplitMetaNamespaceKey(key)
	watcher, err := s.WatcherLister.Watchers(namespace).Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Infof("[watcher] watcher deleted: %s", key)
			if informer, ok := s.WatcherMap[key]; ok {
				informer.ShutDown()
			}
			delete(s.WatcherMap, key)
			return nil
		}
		return errors.Wrap(err, "get watcher error")
	}

	informer, err := s.setupInformer(
		namespace, models.WatcherInformerName(namespace, name), watcher.Spec)
	if err != nil {
		return errors.Wrap(err, "create informer error")
	}
	klog.V(2).Infof("[watcher] watcher added: %s", key)
	go informer.Informer.Run(informer.StopCh)
	cache.WaitForCacheSync(informer.StopCh, informer.Informer.HasSynced)

	for i := 0; i < informer.Workers; i++ {
		go wait.Until(informer.RunWorker, time.Second, informer.StopCh)
	}

	s.WatcherMap[key] = informer
	return nil
}

// handleWatcherErr checks the result, schedules retry if needed
func (s *InformerSet) handleWatcherErr(key string, result error) {
	if result == nil {
		if klog.V(2).Enabled() {
			klog.Infof("[watcher] [%s try] key processed: [%s]",
				humanize.Ordinal(s.WatcherQueue.NumRequeues(key)+1), key)
		}
		// clear retry counter after success
		s.WatcherQueue.Forget(key)
		return
	}

	// retrying
	klog.Warningf("[watcher] [%s try] error processing: [%s]: %s, will be retried",
		humanize.Ordinal(s.WatcherQueue.NumRequeues(key)+1), key, result)
	s.WatcherQueue.AddRateLimited(key)
}
