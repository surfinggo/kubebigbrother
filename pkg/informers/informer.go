package informers

import (
	"fmt"
	"github.com/pkg/errors"
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
	ChannelMap ChannelMap

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
	event := obj.(*Event)
	klog.V(5).Infof("a new item from queue: [%s] %s", event.Type, NamespaceKey(event.Obj))

	i.processingItems.Add(1)

	// we need to mark item as completed whether success or fail
	defer i.Queue.Done(event)

	result := i.processItem(event)
	i.handleErr(event, result)

	i.processingItems.Done()

	return true
}

// processItem process an item synchronously
func (i *Informer) processItem(event *Event) error {
	var channelNamesLeft []ChannelName
	namedErrors := make(map[ChannelName]error)
	for _, channelName := range event.ChannelNames {
		if channel, ok := i.ChannelMap[channelName]; ok {
			if err := channel.Handle(event); err != nil {
				channelNamesLeft = append(channelNamesLeft, channelName)
				namedErrors[channelName] = err
			}
		}
	}

	event.ChannelNames = channelNamesLeft

	// no channels left means process succeeded!
	if len(channelNamesLeft) == 0 {
		return nil
	}
	var s []string
	for channelName, err := range namedErrors {
		s = append(s, fmt.Sprintf("%s: %s", channelName, err))
	}
	return errors.Errorf(strings.Join(s, ","))
}

// handleErr checks the result, schedules retry if needed
func (i *Informer) handleErr(event *Event, result error) {
	if result == nil {
		klog.V(5).Infof("processed: %s", NamespaceKey(event.Obj))
		// clear retry counter after success
		i.Queue.Forget(event)
		return
	}

	if i.Queue.NumRequeues(event) < 3 {
		klog.Warningf("error processing %s: %v", event, result)
		// retrying
		i.Queue.AddRateLimited(event)
		return
	}

	klog.Error(fmt.Errorf("max retries exceeded, "+
		"dropping item %s out of the queue: %v", event, result))
	// max retries exceeded, forget it
	i.Queue.Forget(event)
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
