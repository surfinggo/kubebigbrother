package controller

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
)

type Config struct {
	DBDialect       string
	DBArgs          string
	KubeConfig      string
	InformersConfig *informers.ConfigFile
}

type Controller struct {
	EventStore event_store.Interface
	Informers  informers.Interface
}

func (c *Controller) Start(stopCh <-chan struct{}) error {
	return c.Informers.Start(stopCh)
}

func (c *Controller) Shutdown() {
	c.Informers.Shutdown()
}

func Setup(config Config) (*Controller, error) {
	recorder := &Controller{}

	db, err := gormdb.New(config.DBDialect, config.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	recorder.EventStore = event_store.New(db)

	informerInstance, err := informers.Setup(informers.Config{
		KubeConfig: config.KubeConfig,
		ConfigFile: config.InformersConfig,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}
	recorder.Informers = informerInstance

	return recorder, nil
}
