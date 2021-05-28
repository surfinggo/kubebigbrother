package channels

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"k8s.io/klog/v2"
)

// ChannelGroup is a special type of channel, which includes a slice of channels
type ChannelGroup struct {
	Channels []Channel
}

// Handle implements Channel
func (c *ChannelGroup) Handle(e *event.Event) error {
	for _, channel := range c.Channels {
		err := channel.Handle(e)
		if err != nil {
			klog.Warning(errors.Wrap(err, "handle event error"))
		}
	}
	// TODO: if an error is returned, all channels in c.ChannelNames will be retried.
	// need to find a way to handle error but only retry failed channels.
	return nil
}

func NewChannelGroup() (*ChannelGroup, error) {
	return &ChannelGroup{
		Channels: nil, // TODO: set channels
	}, nil
}
