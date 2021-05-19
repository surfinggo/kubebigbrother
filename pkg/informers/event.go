package informers

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type EventType string

const (
	EventTypeAdded   = "ADDED"
	EventTypeDeleted = "DELETED"
	EventTypeUpdated = "UPDATED"
)

type Event struct {
	// Type is the type of the event
	Type EventType

	// Obj is the resource affected
	Obj *unstructured.Unstructured

	// OldObj is only set for EventTypeUpdated
	OldObj *unstructured.Unstructured
}
