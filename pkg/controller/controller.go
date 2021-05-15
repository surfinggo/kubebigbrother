package controller

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
)

type Options struct {
	DBDialect       string
	DBArgs          string
	KubeConfig      string
	InformersConfig *informers.Config
}

type Controller struct {
	EventStore event_store.Interface
	Informers  informers.Interface
}

func (r *Controller) Start(stopCh <-chan struct{}) error {
	return r.Informers.Start(stopCh)
}

func Setup(options Options) (*Controller, error) {
	recorder := &Controller{}

	db, err := gormdb.New(options.DBDialect, options.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	recorder.EventStore = event_store.New(db)

	informerInstance, err := informers.Setup(informers.Options{
		KubeConfig: options.KubeConfig,
		Config:     options.InformersConfig,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}
	recorder.Informers = informerInstance

	return recorder, nil
}
