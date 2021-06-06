package informers

import (
	"github.com/pkg/errors"
	spg "github.com/spongeprojects/client-go/api/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"reflect"
	"strings"
	"sync"
)

func (s *InformerSet) setupInformer(
	namespace, informerName string, c spg.WatcherSpec) (*Informer, error) {
	channelNames := c.ChannelNames
	if len(channelNames) == 0 {
		channelNames = s.DefaultChannelNames
	}

	resyncPeriod, set, err := parseResyncPeriod(c.ResyncPeriod)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid resync period: %s", c.ResyncPeriod)
	}
	if !set {
		resyncPeriod = s.DefaultResyncPeriodFunc()
	}

	workers := c.Workers
	if workers < 1 {
		workers = s.DefaultWorkers
	}

	maxRetries := c.MaxRetries
	if maxRetries < 1 {
		maxRetries = s.DefaultMaxRetries
	}

	gvr, err := s.ResourceBuilder.ParseGroupResource(c.Resource)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid resource: %s", c.Resource)
	}

	saveSilently := func(e *event.Event) {
		s.EventStore.SaveSilently(e.ToModel(informerName, gvr))
	}

	rateLimiter := workqueue.DefaultControllerRateLimiter()
	queue := workqueue.NewRateLimitingQueue(rateLimiter)
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		s.DynamicClient, resyncPeriod, namespace, nil)
	resourceInformer := informerFactory.ForResource(gvr).Informer()

	handlerFuncs := cache.ResourceEventHandlerFuncs{}
	if !c.NoticeWhenAdded &&
		!c.NoticeWhenDeleted &&
		!c.NoticeWhenUpdated {
		// for some reason, informer won't start if they are all false,
		// maybe someday someone can clarify why this happen.
		return nil, errors.Errorf(
			"NoticeWhenAdded, NoticeWhenDeleted and NoticeWhenUpdated " +
				"cannot be false simultaneously")
	}
	if c.NoticeWhenAdded {
		handlerFuncs.AddFunc = func(obj interface{}) {
			st, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			e := event.NewAdded(st)

			if !s.JustWatch {
				isCurrentlyAdded, err := s.EventStore.IsCurrentlyAdded(
					informerName, gvr.Group, gvr.Version, gvr.Resource,
					e.Obj.GetNamespace(), e.Obj.GetName())
				if err != nil {
					klog.Warning(errors.Wrap(err, "find latest record error"))
				} else if isCurrentlyAdded {
					klog.V(5).Infof(
						"[%s] resource is currently added, skip ADDED event: [%s] [%s]",
						informerName, e.Type, e.GroupVersionKindName())
					return // ADDED event are emitted when controller restart
				}
			}

			klog.V(5).Infof("[%s] received: [%s] [%s]",
				informerName, e.Type, e.GroupVersionKindName())

			if !s.JustWatch {
				saveSilently(e)
			}

			queue.Add(s.wrap(e, channelNames))
		}
	}
	if c.NoticeWhenDeleted {
		handlerFuncs.DeleteFunc = func(obj interface{}) {
			st, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			e := event.NewDeleted(st)

			klog.V(5).Infof("[%s] received: [%s] [%s]",
				informerName, e.Type, utils.GroupVersionKindName(st))

			if !s.JustWatch {
				saveSilently(e)
			}

			queue.Add(s.wrap(e, channelNames))
		}
	}
	if c.NoticeWhenUpdated {
		handlerFuncs.UpdateFunc = func(oldObj, newObj interface{}) {
			oldSt, ok1 := oldObj.(*unstructured.Unstructured)
			st, ok2 := newObj.(*unstructured.Unstructured)
			if !ok1 || !ok2 {
				return
			}
			updated := false
			for _, field := range c.UpdateOn {
				fieldPath := strings.Split(strings.TrimPrefix(field, "."), ".")
				f1, exist1, err1 := unstructured.NestedFieldNoCopy(
					st.Object, fieldPath...)
				f2, exist2, err2 := unstructured.NestedFieldNoCopy(
					oldSt.Object, fieldPath...)
				if !exist1 || !exist2 {
					klog.Warningf("[%s] field not exist in resource: %s: %s",
						informerName, utils.GroupVersionKindName(st), field)
				}
				if err1 != nil {
					klog.Warningf("[%s] get field value error, resource: %s: %s",
						informerName, utils.GroupVersionKindName(st), err1)
				}
				if err2 != nil {
					klog.Warningf("[%s] get field value error, resource: %s: %s",
						informerName, utils.GroupVersionKindName(st), err2)
				}
				if exist1 && exist2 && err1 == nil && err2 == nil &&
					!reflect.DeepEqual(f1, f2) {
					updated = true
					break
				}
			}
			if c.UpdateOn == nil || updated {
				e := event.NewUpdated(st, oldSt)

				klog.V(5).Infof("[%s] received: [%s] [%s]",
					informerName, e.Type, utils.GroupVersionKindName(st))

				if !s.JustWatch {
					saveSilently(e)
				}

				queue.Add(s.wrap(e, channelNames))
			}
		}
	}
	resourceInformer.AddEventHandler(handlerFuncs)

	return &Informer{
		ID:              informerName,
		Resource:        c.Resource,
		GVR:             gvr,
		UpdateOn:        c.UpdateOn,
		ChannelMap:      s.ChannelMap,
		Informer:        resourceInformer,
		Queue:           queue,
		Workers:         workers,
		MaxRetries:      maxRetries,
		processingItems: &sync.WaitGroup{}, // TODO: wait before exit
		StopCh:          make(chan struct{}),
	}, nil
}
