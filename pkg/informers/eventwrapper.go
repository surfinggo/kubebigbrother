package informers

import (
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

type eventWrapper struct {
	*event.Event

	// ChannelNames is channels to process,
	// name will be removed from slice after processed successfully
	ChannelNames []channels.ChannelName
}
