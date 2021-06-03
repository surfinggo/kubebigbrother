package models

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"time"
)

// Event is a copy of Kubernetes event, to persistent event history,
// and to avoid improper ADDED events when the controller restart.
type Event struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`

	// InformerName is unique value represents the informer config,
	// every Event belongs to an informer.
	InformerName string `json:"informer_name"`

	EventType string `json:"event_type"`
	Group     string `gorm:"column:event_group" json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Obj       []byte `json:"obj,omitempty"`
	OldObj    []byte `json:"old_obj,omitempty"`
}

func (e *Event) GetObj() (obj *unstructured.Unstructured) {
	if e.Obj == nil {
		return nil
	}
	_ = json.Unmarshal(e.Obj, &obj)
	return obj
}

func (e *Event) GetOldObj() (obj *unstructured.Unstructured) {
	if e.OldObj == nil {
		return nil
	}
	_ = json.Unmarshal(e.OldObj, &obj)
	return obj
}
