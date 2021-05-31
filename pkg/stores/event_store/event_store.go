package event_store

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

type Interface interface {
	List() (events []models.Event, err error)
	Save(event *models.Event) (err error)
	SaveSilently(event *models.Event)
}

type Store struct {
	DB *gorm.DB
}

func (s *Store) List() (events []models.Event, err error) {
	err = s.DB.Order("id desc").Limit(20).Find(&events).Error
	return
}

func (s *Store) Save(event *models.Event) error {
	return s.DB.Save(event).Error
}

func (s *Store) SaveSilently(event *models.Event) {
	if err := s.Save(event); err != nil {
		klog.Warning(errors.Wrap(err, "save event error"))
	}
}

func New(db *gorm.DB) Interface {
	return &Store{
		DB: db,
	}
}
