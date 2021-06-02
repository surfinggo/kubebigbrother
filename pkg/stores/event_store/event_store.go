package event_store

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

type Interface interface {
	List() (events []models.Event, err error)
	ListByInformer(informerName string) (events []models.Event, err error)
	Search(q string) (events []models.Event, err error)
	IsCurrentlyAdded(informerName,
		group, version, resource, namespace, name string) (exist bool, err error)
	Save(event *models.Event) (err error)
	SaveSilently(event *models.Event)
}

type Store struct {
	DB *gorm.DB
}

func (s *Store) List() (events []models.Event, err error) {
	err = s.DB.Order("id desc").Limit(50).Find(&events).Error
	return
}

func (s *Store) ListByInformer(informerName string) (events []models.Event, err error) {
	err = s.DB.Where("informer_name = ?", informerName).Order("id desc").Limit(50).Find(&events).Error
	return
}

func (s *Store) Search(q string) (events []models.Event, err error) {
	l := fmt.Sprintf("%%%s%%", q)
	err = s.DB.Where("name like ?", l).
		Or("namespace like ?", l).
		Or("event_group like ?", l).
		Or("version like ?", l).
		Or("resource like ?", l).
		Order("id desc").Limit(50).Find(&events).Error
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

func (s *Store) IsCurrentlyAdded(informerName,
	group, version, resource, namespace, name string) (yes bool, err error) {
	var e models.Event
	if err := s.DB.Where("informer_name = ?", informerName).
		// TODO: use event.EventType ADDED and DELETED without import loop
		Where("event_type in ?", []string{"ADDED", "DELETED"}).
		Where("event_group = ?", group).
		Where("version = ?", version).
		Where("resource = ?", resource).
		Where("namespace = ?", namespace).
		Where("name = ?", name).
		Order("create_time desc").
		First(&e).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	// the latest event type is ADDED means the kind is in cache
	return e.EventType == "ADDED", nil
}

func New(db *gorm.DB) Interface {
	return &Store{
		DB: db,
	}
}
