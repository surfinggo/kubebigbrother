package watcher

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"os"
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

func Setup(options Options) (*Watcher, error) {
	watcher := &Watcher{}

	config := options.InformersConfig

	channelNamePrint := informers.ChannelName("print")

	adapted := &informers.Config{
		Channels: map[informers.ChannelName]informers.ChannelConfig{
			channelNamePrint: {
				Name: channelNamePrint,
				Type: informers.ChannelTypePrint,
				Print: &informers.ChannelPrintConfig{
					Out: os.Stdout,
				},
			},
		},
		MinResyncPeriod: config.MinResyncPeriod,
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
			})
		}
		adapted.Namespaces = append(adapted.Namespaces, informers.NamespaceConfig{
			Namespace: namespace.Namespace,
			Resources: resources,
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
