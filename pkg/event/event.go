package event

import (
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Type string

const (
	TypeAdded   = "ADDED"
	TypeDeleted = "DELETED"
	TypeUpdated = "UPDATED"
)

type Event struct {
	// Type is the type of the event
	Type Type

	// Obj is the resource affected
	Obj *unstructured.Unstructured

	// OldObj is only set for EventTypeUpdated
	OldObj *unstructured.Unstructured

	// gvkNameCache is a cache for GroupVersionKindName
	gvkNameCache string
}

// GroupVersionKindName returns group, version, kind, namespace and name string
// for the affected resource
//
// examples:
//   /v1, Kind=ConfigMap, demo/demo
//   apps/v1, Kind=Deployment, demo/canary
func (e *Event) GroupVersionKindName() string {
	if e.gvkNameCache == "" {
		e.gvkNameCache = utils.GroupVersionKindName(e.Obj)
	}
	return e.gvkNameCache
}

// NamespaceKey returns namespaced key for the affected resource
func (e *Event) NamespaceKey() string {
	return utils.NamespaceKey(e.Obj)
}

func NewAdded(obj *unstructured.Unstructured) *Event {
	return &Event{
		Type: TypeAdded,
		Obj:  obj,
	}
}

func NewDeleted(obj *unstructured.Unstructured) *Event {
	return &Event{
		Type: TypeDeleted,
		Obj:  obj,
	}
}

func NewUpdated(obj, oldObj *unstructured.Unstructured) *Event {
	return &Event{
		Type:   TypeUpdated,
		Obj:    obj,
		OldObj: oldObj,
	}
}
