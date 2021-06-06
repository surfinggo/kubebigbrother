package controller

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"time"
)

type Config struct {
	DBDialect           string
	DBArgs              string
	Kubeconfig          string
	DefaultWorkers      int
	DefaultMaxRetries   int
	DefaultChannelNames []string
	MinResyncPeriod     time.Duration
}

type Controller struct {
	EventStore event_store.Interface
	Informers  informers.Interface
}

// Start starts the Controller
func (c *Controller) Start(stopCh <-chan struct{}) error {
	return c.Informers.Start(stopCh)
}

// Shutdown shutdowns the Controller
func (c *Controller) Shutdown() {
	c.Informers.Shutdown()
}

// Setup sets up a new Controller
func Setup(config Config) (*Controller, error) {
	controller := &Controller{}

	db, err := gormdb.New(config.DBDialect, config.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	controller.EventStore = event_store.New(db)

	informerInstance, err := informers.Setup(informers.Config{
		Kubeconfig:          config.Kubeconfig,
		DefaultWorkers:      config.DefaultWorkers,
		DefaultMaxRetries:   config.DefaultMaxRetries,
		DefaultChannelNames: config.DefaultChannelNames,
		MinResyncPeriod:     config.MinResyncPeriod,
		JustWatch:           false,
		EventStore:          controller.EventStore,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}
	controller.Informers = informerInstance

	return controller, nil
}
