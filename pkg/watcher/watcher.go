package watcher

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
)

type Options struct {
	KubeConfig      string
	InformersConfig *informers.Config
}

type Watcher struct {
	Informers informers.Interface
}

func (w *Watcher) Start(stopCh <-chan struct{}) {
	w.Informers.Start(stopCh)
}

func Setup(options Options) (*Watcher, error) {
	watcher := &Watcher{}

	informerSet, err := informers.Setup(informers.Options{
		KubeConfig: options.KubeConfig,
		Config:     options.InformersConfig,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}

	watcher.Informers = informerSet

	return watcher, nil
}
