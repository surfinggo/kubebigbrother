package watcher

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
)

type Config struct {
	Kubeconfig string
}

type Watcher struct {
	Informers informers.Interface
}

func (w *Watcher) Start(stopCh <-chan struct{}) error {
	return w.Informers.Start(stopCh)
}

func (w *Watcher) Shutdown() {
	w.Informers.Shutdown()
}

func Setup(config Config) (*Watcher, error) {
	watcher := &Watcher{}

	informerSet, err := informers.Setup(informers.Config{
		Kubeconfig: config.Kubeconfig,
		JustWatch:  true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}

	watcher.Informers = informerSet

	return watcher, nil
}
