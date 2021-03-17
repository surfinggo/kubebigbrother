package watcher

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Version string

	Env string

	GinDebug   bool
	KubeConfig string
	Resource   string
}

type Watcher struct {
	Version string

	Addr string

	Clientset kubernetes.Interface

	Controller interface{}
}

func SetupWatcher(options Options) (*Watcher, error) {
	app := &Watcher{}
	app.Addr = options.Resource
	app.Version = options.Version
	app.Controller = nil

	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "create kube client error")
	}

	app.Clientset = clientset

	return app, nil
}
