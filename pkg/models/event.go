package models

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"time"
)

type Event struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	EventType event.Type `json:"event_type"`
	Group     string     `json:"group"`
	Version   string     `json:"version"`
	Kind      string     `json:"kind"`
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Obj       []byte     `json:"obj"`
	OldObj    []byte     `json:"old_obj"`
}
