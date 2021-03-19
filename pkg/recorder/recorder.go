package recorder

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/informer"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	v1 "k8s.io/api/core/v1"
)

type Options struct {
	Env string

	DBDialect  string
	DBArgs     string
	KubeConfig string
	Resource   string
}

type Recorder struct {
	EventStore event_store.Interface
	Informer   *informer.Informer
}

func Setup(options Options) (*Recorder, error) {
	recorder := &Recorder{}

	db, err := gormdb.New(options.DBDialect, options.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	recorder.EventStore = event_store.New(db)

	informerInstance, err := informer.Setup(informer.Options{
		KubeConfig: options.KubeConfig,
		Resource:   options.Resource,
		ConfigMapAddFunc: func(configMap *v1.ConfigMap) {
			log.Infof("created: %s/%s", configMap.Namespace, configMap.Name)
			e := &models.Event{
				Description: fmt.Sprintf("created %s/%s", configMap.Namespace, configMap.Name),
			}
			err := recorder.EventStore.Save(e)
			if err != nil {
				log.Error(errors.Wrap(err, "save event error"))
			}
		},
		ConfigMapUpdateFunc: func(oldConfigMap *v1.ConfigMap, newConfigMap *v1.ConfigMap) {
			log.Infof("updated: %s/%s", newConfigMap.Namespace, newConfigMap.Name)
			e := &models.Event{
				Description: fmt.Sprintf("updated %s/%s", newConfigMap.Namespace, newConfigMap.Name),
			}
			err := recorder.EventStore.Save(e)
			if err != nil {
				log.Error(errors.Wrap(err, "save event error"))
			}
		},
		ConfigMapDeleteFunc: func(configMap *v1.ConfigMap) {
			log.Infof("deleted: %s/%s", configMap.Namespace, configMap.Name)
			e := &models.Event{
				Description: fmt.Sprintf("deleted %s/%s", configMap.Namespace, configMap.Name),
			}
			err := recorder.EventStore.Save(e)
			if err != nil {
				log.Error(errors.Wrap(err, "save event error"))
			}
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informer error")
	}
	recorder.Informer = informerInstance

	return recorder, nil
}

func (r *Recorder) Start() error {
	return r.Informer.Start()
}
