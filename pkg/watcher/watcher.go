package watcher

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
)

const (
	channelNamePrint = informers.ChannelName("print")
)

type Options struct {
	KubeConfig      string
	InformersConfig *informers.Config
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

func Setup(options Options) (*Watcher, error) {
	watcher := &Watcher{}

	config := options.InformersConfig

	// watcher only prints event to stdout, all user defined channels will be
	// dropped.
	adapted := &informers.Config{
		Channels: map[informers.ChannelName]informers.ChannelConfig{
			channelNamePrint: {
				Type: informers.ChannelTypePrint,
				Print: &informers.ChannelPrintConfig{
					Writer: informers.PrintWriterStdout,
				},
			},
		},
		DefaultWorkers:      config.DefaultWorkers,
		DefaultChannelNames: config.DefaultChannelNames,
		MinResyncPeriod:     config.MinResyncPeriod,
	}
	for _, namespace := range config.Namespaces {
		var resources []informers.ResourceConfig
		for _, resource := range namespace.Resources {
			resources = append(resources, informers.ResourceConfig{
				Resource:          resource.Resource,
				NoticeWhenAdded:   resource.NoticeWhenAdded,
				NoticeWhenDeleted: resource.NoticeWhenDeleted,
				NoticeWhenUpdated: resource.NoticeWhenUpdated,
				UpdateOn:          resource.UpdateOn,
				ChannelNames:      []informers.ChannelName{channelNamePrint},
				ResyncPeriod:      resource.ResyncPeriod,
				Workers:           resource.Workers,
			})
		}
		adapted.Namespaces = append(adapted.Namespaces, informers.NamespaceConfig{
			Namespace:           namespace.Namespace,
			Resources:           resources,
			DefaultWorkers:      namespace.DefaultWorkers,
			DefaultChannelNames: namespace.DefaultChannelNames,
			MinResyncPeriod:     namespace.MinResyncPeriod,
		})
	}

	informerSet, err := informers.Setup(informers.Options{
		KubeConfig: options.KubeConfig,
		Config:     adapted,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}

	watcher.Informers = informerSet

	return watcher, nil
}
