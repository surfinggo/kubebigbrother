package models

import (
	"time"
)

// Event is a copy of Kubernetes event, to persistent event history,
// and to avoid improper ADDED events when the controller restart.
type Event struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`

	// InformerConfigHash is unique value represents the informer config,
	// every Event belongs to an informer.
	InformerConfigHash string `json:"informer_config_hash"`

	EventType string `json:"event_type"`
	Group     string `gorm:"column:event_group" json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Obj       []byte `json:"obj"`
	OldObj    []byte `json:"old_obj"`
}
