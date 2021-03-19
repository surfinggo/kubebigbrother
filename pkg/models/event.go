package models

import (
	"time"
)

type Event struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Description string `json:"description"`
}
