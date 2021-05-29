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

type Config struct {
	KubeConfig      string
	InformersConfig *informers.ConfigFile
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
	if len(config.InformersConfig.Channels) != 0 {
		var p string
		if len(config.InformersConfig.Channels) == 1 {
			p = "the channel has"
		} else {
			p = fmt.Sprintf("%d channels have", len(config.InformersConfig.Channels))
		}
		klog.Warningf("watch: %s been replaced by %s", p, channelNamePrintToStdout)
	}

	watcher := &Watcher{}

	informersConfig := config.InformersConfig

	// watcher only prints event to stdout, all user defined channels will be
	// dropped.
	adapted := &informers.ConfigFile{
		Channels: map[channels.ChannelName]informers.ChannelConfig{
			channelNamePrintToStdout: {
				Type: channels.ChannelTypePrint,
				Print: &informers.ChannelPrintConfig{
					Writer: channels.PrintWriterStdout,
				},
			},
		},
		DefaultWorkers:      informersConfig.DefaultWorkers,
		DefaultChannelNames: informersConfig.DefaultChannelNames,
		MinResyncPeriod:     informersConfig.MinResyncPeriod,
	}
	for _, namespace := range informersConfig.Namespaces {
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

	informerSet, err := informers.Setup(informers.Config{
		KubeConfig: config.KubeConfig,
		ConfigFile: adapted,
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup informers error")
	}

	watcher.Informers = informerSet

	return watcher, nil
}
