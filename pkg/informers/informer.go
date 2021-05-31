package informers

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"strings"
	"sync"
)

type Informer struct {
	// ID is an unique string to identify instances
	ID string

	// Resource is the resource to watch, e.g. "deployments.v1.apps"
	Resource string

	// GVR is group version kind of Resource
	GVR schema.GroupVersionResource

	// UpdateOn defines fields to watch, used with NoticeWhenUpdated
	UpdateOn []string

	// ChannelMap defines channels to send notification
	ChannelMap channels.ChannelMap

	// Queue is a rate limiting queue
	Queue workqueue.RateLimitingInterface

	processingItems *sync.WaitGroup

	// Workers is number of workers
	Workers int

	// MaxRetries defines max retry times
	MaxRetries int
}

func (i *Informer) RunWorker() {
	for i.processNextItem() {
	}
}

// processNextItem waits and processes items in the queue
func (i *Informer) processNextItem() bool {
	// block until an item arrives or queue shutdown
	obj, shutdown := i.Queue.Get()
	if shutdown {
		return false
	}
	item := obj.(*eventWrapper)

	if klog.V(5).Enabled() {
		klog.Infof("[%s] [%s try] item pop from queue: [%s] [%s]",
			i.ID, humanize.Ordinal(i.Queue.NumRequeues(item)+1),
			item.Event.Type, item.GroupVersionKindName())
	}

	i.processingItems.Add(1)

	// we need to mark item as completed whether success or fail
	defer i.Queue.Done(item)

	result := i.processItem(item)
	i.handleErr(item, result)

	i.processingItems.Done()

	return true
}

// processItem process an item synchronously
func (i *Informer) processItem(item *eventWrapper) error {
	errs := make(map[ChannelToProcess]error)
	for _, ch := range item.ChannelsToProcess {
		if channel, ok := i.ChannelMap[ch.ChannelName]; ok {
			if err := channel.Handle(ch.EventProcessContext); err != nil {
				errs[ch] = err
			}
		}
	}

	if len(errs) == 0 { // no error, no channel left, everything works as expected
		item.ChannelsToProcess = nil
		return nil
	}

	var channelToProcessLeft []ChannelToProcess
	var es []string
	for ch, err := range errs {
		channelToProcessLeft = append(channelToProcessLeft, ch)
		es = append(es, fmt.Sprintf("channel %s error: %s", ch.ChannelName, err))
	}
	item.ChannelsToProcess = channelToProcessLeft
	return errors.Errorf("process error: %s", strings.Join(es, ","))
}

// handleErr checks the result, schedules retry if needed
func (i *Informer) handleErr(item *eventWrapper, result error) {
	if result == nil {
		if klog.V(5).Enabled() {
			klog.Infof("[%s] [%s try] item processed: [%s] [%s]",
				i.ID, humanize.Ordinal(i.Queue.NumRequeues(item)+1),
				item.Event.Type, item.GroupVersionKindName())
		}
		// clear retry counter after success
		i.Queue.Forget(item)
		return
	}

	if i.Queue.NumRequeues(item) >= i.MaxRetries-1 {
		klog.Errorf(
			"[%s] [%s try] error processing: "+
				"[%s] [%s]: %s, max retries exceeded, dropping item out of the queue",
			i.ID, humanize.Ordinal(i.Queue.NumRequeues(item)+1),
			item.Event.Type, item.GroupVersionKindName(), result)

		// max retries exceeded, forget it
		i.Queue.Forget(item)
		return
	}

	if klog.V(5).Enabled() {
		klog.Warningf("[%s] [%s try] error processing: [%s] [%s]: %s, will be retried",
			i.ID, humanize.Ordinal(i.Queue.NumRequeues(item)+1),
			item.Event.Type, item.GroupVersionKindName(), result)
	}
	// retrying
	i.Queue.AddRateLimited(item)
}

func (i *Informer) ShutDown() {
	i.Queue.ShutDown()
}
