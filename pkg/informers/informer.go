package informers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"strings"
	"sync"
)

type InformerInterface interface {
	RunWorker()
	GetWorkers() int
	GetResource() string
	ShutDown()
}

type Informer struct {
	// ID is an unique string to identify instances
	ID string

	// Resource is the resource to watch, e.g. "deployments.v1.apps"
	Resource string

	// UpdateOn defines fields to watch, used with NoticeWhenUpdated
	UpdateOn []string

	// ChannelMap defines channels to send notification
	ChannelMap channels.ChannelMap

	// Queue is a rate limiting queue
	Queue workqueue.RateLimitingInterface

	processingItems *sync.WaitGroup

	// Workers is number of workers
	Workers int
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
	klog.V(5).Infof("new item from queue: [%s] %s", item.Event.Type, item.GroupVersionKindName())

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
	for _, channelToProcess := range item.ChannelsToProcess {
		if channel, ok := i.ChannelMap[channelToProcess.ChannelName]; ok {
			if err := channel.Handle(channelToProcess.EventProcessContext); err != nil {
				errs[channelToProcess] = err
			}
		}
	}

	if len(errs) == 0 { // no error, no channel left, everything works as expected
		item.ChannelsToProcess = nil
		return nil
	}

	var channelToProcessLeft []ChannelToProcess
	var es []string
	for channelToProcess, err := range errs {
		channelToProcessLeft = append(channelToProcessLeft, channelToProcess)
		es = append(es, fmt.Sprintf("channel %s error: %s", channelToProcess.ChannelName, err))
	}
	item.ChannelsToProcess = channelToProcessLeft
	return errors.Errorf("process error: %s", strings.Join(es, ","))
}

// handleErr checks the result, schedules retry if needed
func (i *Informer) handleErr(item *eventWrapper, result error) {
	if result == nil {
		klog.V(5).Infof("processed: [%s] %s", item.Event.Type, item.GroupVersionKindName())
		// clear retry counter after success
		i.Queue.Forget(item)
		return
	}

	if i.Queue.NumRequeues(item) <= 3 {
		klog.Warningf("error processing [%s] %s: %v",
			item.Event.Type, item.GroupVersionKindName(), result)
		// retrying
		i.Queue.AddRateLimited(item)
		return
	}

	klog.Error(fmt.Errorf(
		"max retries exceeded, dropping item [%s] out of the queue: %v",
		item.GroupVersionKindName(), result))
	// max retries exceeded, forget it
	i.Queue.Forget(item)
}

func (i *Informer) GetWorkers() int {
	return i.Workers
}

func (i *Informer) GetResource() string {
	return i.Resource
}

func (i *Informer) ShutDown() {
	i.Queue.ShutDown()
}
