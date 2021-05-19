package watcher

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"k8s.io/klog/v2"
)

const (
	channelNamePrintToStdout = channels.ChannelName("print-to-stdout")
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
	if len(options.InformersConfig.Channels) != 0 {
		var p string
		if len(options.InformersConfig.Channels) == 1 {
			p = "the channel has"
		} else {
			p = fmt.Sprintf("%d channels have", len(options.InformersConfig.Channels))
		}
		klog.Warningf("watch: %s been replaced by a single channel: %s", p, channelNamePrintToStdout)
	}

	watcher := &Watcher{}

	config := options.InformersConfig

	// watcher only prints event to stdout, all user defined channels will be
	// dropped.
	adapted := &informers.Config{
		Channels: map[channels.ChannelName]informers.ChannelConfig{
			channelNamePrintToStdout: {
				Type: channels.ChannelTypePrint,
				Print: &informers.ChannelPrintConfig{
					Writer: channels.PrintWriterStdout,
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
				ChannelNames:      []channels.ChannelName{channelNamePrintToStdout},
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
