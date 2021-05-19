package event

import (
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
