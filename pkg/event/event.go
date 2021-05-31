package event

import (
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"github.com/spongeprojects/kubebigbrother/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Type of event
type Type string

const (
	TypeAdded   = "ADDED"   // resource created
	TypeDeleted = "DELETED" // resource deleted
	TypeUpdated = "UPDATED" // resource updated
)

// Event is representation of Kubernetes event
type Event struct {
	// Type is the type of the event
	Type Type `json:"type"`

	// Obj is the resource affected
	Obj *unstructured.Unstructured `json:"obj"`

	// OldObj is only set for EventTypeUpdated
	OldObj *unstructured.Unstructured `json:"oldObj,omitempty"`

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

// Color returns theme color for the type of the event
func (e *Event) Color() string {
	switch e.Type {
	case TypeAdded:
		return style.Success
	case TypeDeleted:
		return style.Warning
	default:
		return style.Info
	}
}

// ToModel translate Event into *models.Event
func (e *Event) ToModel(informerConfigHash string) *models.Event {
	model := &models.Event{
		InformerConfigHash: informerConfigHash,
		EventType:          string(e.Type),
		Namespace:          e.Obj.GetNamespace(),
		Name:               e.Obj.GetName(),
	}

	gvk := e.Obj.GroupVersionKind()
	model.Group = gvk.Group
	model.Version = gvk.Version
	model.Kind = gvk.Kind

	objJSON, _ := e.Obj.MarshalJSON()
	model.Obj = objJSON

	if e.OldObj != nil {
		oldObjJSON, _ := e.OldObj.MarshalJSON()
		model.OldObj = oldObjJSON
	}

	return model
}

// NewAdded creates added event
func NewAdded(obj *unstructured.Unstructured) *Event {
	return &Event{
		Type: TypeAdded,
		Obj:  obj,
	}
}

// NewDeleted creates deleted event
func NewDeleted(obj *unstructured.Unstructured) *Event {
	return &Event{
		Type: TypeDeleted,
		Obj:  obj,
	}
}

// NewUpdated creates updated event
func NewUpdated(obj, oldObj *unstructured.Unstructured) *Event {
	return &Event{
		Type:   TypeUpdated,
		Obj:    obj,
		OldObj: oldObj,
	}
}
