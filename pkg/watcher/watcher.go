package watcher

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informer"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	v1 "k8s.io/api/core/v1"
)

type Options struct {
	Env string

	KubeConfig string
	Resource   string
}

type Watcher struct {
	Informer *informer.Informer
}

func Setup(options Options) (*Watcher, error) {
	watcher := &Watcher{}

	informerInstance, err := informer.Setup(informer.Options{
		KubeConfig: options.KubeConfig,
		Resource:   options.Resource,
		ConfigMapAddFunc: func(configMap *v1.ConfigMap) {
			log.Infof("created: %s/%s", configMap.Namespace, configMap.Name)
		},
		ConfigMapUpdateFunc: func(oldConfigMap *v1.ConfigMap, newConfigMap *v1.ConfigMap) {
			log.Infof("updated: %s/%s", newConfigMap.Namespace, newConfigMap.Name)
		},
		ConfigMapDeleteFunc: func(configMap *v1.ConfigMap) {
			log.Infof("deleted: %s/%s", configMap.Namespace, configMap.Name)
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informer error")
	}
	watcher.Informer = informerInstance

	return watcher, nil
}

func (w *Watcher) Start() error {
	return w.Informer.Start()
}
