package event_store

import (
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	List() (events []models.Event, err error)
	Save(event *models.Event) (err error)
}

type Store struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Interface {
	return &Store{
		DB: db,
	}
}

func (s *Store) List() (events []models.Event, err error) {
	err = s.DB.Order("id desc").Limit(20).Find(&events).Error
	return
}

func (s *Store) Save(event *models.Event) error {
	return s.DB.Save(event).Error
}
