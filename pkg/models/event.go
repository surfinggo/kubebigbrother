package models

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"time"
)

// Event is a copy of Kubernetes event, to persistent event history,
// and to avoid improper ADDED events when the controller restart.
type Event struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	// InformerChecksum is md5 value of an informer config,
	// every Event belongs to an informer.
	InformerChecksum string `json:"informer_checksum"`

	EventType event.Type `json:"event_type"`
	Group     string     `json:"group"`
	Version   string     `json:"version"`
	Kind      string     `json:"kind"`
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Obj       []byte     `json:"obj"`
	OldObj    []byte     `json:"old_obj"`
}
