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

func (s *InformerSet) RunClusterWatcherWorker() {
	for s.processNextClusterWatcher() {
	}
}

// processNextClusterWatcher waits and processes items in the queue
func (s *InformerSet) processNextClusterWatcher() bool {
	// block until an key arrives or queue shutdown
	obj, shutdown := s.ClusterWatcherQueue.Get()
	if shutdown {
		return false
	}
	key := obj.(string)
	if klog.V(2).Enabled() {
		klog.Infof("[clusterwatcher] [%s try] key pop from queue: [%s]",
			humanize.Ordinal(s.ClusterWatcherQueue.NumRequeues(key)+1), key)
	}

	// we need to mark key as completed whether success or fail
	defer s.ClusterWatcherQueue.Done(key)

	result := s.processClusterWatcher(key)
	s.handleClusterWatcherErr(key, result)

	return true
}

// processClusterWatcher process an item synchronously
func (s *InformerSet) processClusterWatcher(key string) error {
	watcher, err := s.ClusterWatcherLister.Get(key)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Infof("[clusterwatcher] clusterwatcher deleted: %s", key)
			if informer, ok := s.ClusterWatcherMap[key]; ok {
				informer.ShutDown()
			}
			delete(s.ClusterWatcherMap, key)
			return nil
		}
		return errors.Wrap(err, "get watcher error")
	}

	informer, err := s.setupInformer(
		"", models.ClusterWatcherInformerName(key), watcher.Spec)
	if err != nil {
		return errors.Wrap(err, "create informer error")
	}
	klog.V(2).Infof("[clusterwatcher] clusterwatcher added: %s", key)
	go informer.Informer.Run(informer.StopCh)
	cache.WaitForCacheSync(informer.StopCh, informer.Informer.HasSynced)

	for i := 0; i < informer.Workers; i++ {
		go wait.Until(informer.RunWorker, time.Second, informer.StopCh)
	}

	s.ClusterWatcherMap[key] = informer
	return nil
}

// handleClusterWatcherErr checks the result, schedules retry if needed
func (s *InformerSet) handleClusterWatcherErr(key string, result error) {
	if result == nil {
		if klog.V(5).Enabled() {
			klog.Infof("[clusterwatcher] [%s try] key processed: [%s]",
				humanize.Ordinal(s.ClusterWatcherQueue.NumRequeues(key)+1), key)
		}
		// clear retry counter after success
		s.ClusterWatcherQueue.Forget(key)
		return
	}

	// retrying
	klog.Warningf("[clusterwatcher] [%s try] error processing: [%s]: %s, will be retried",
		humanize.Ordinal(s.ClusterWatcherQueue.NumRequeues(key)+1), key, result)
	s.ClusterWatcherQueue.AddRateLimited(key)
}
