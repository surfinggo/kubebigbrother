package informers

import (
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelToProcess defines a channel to process
type ChannelToProcess struct {
	ChannelName         ChannelName
	EventProcessContext *channels.EventProcessContext
}

type eventWrapper struct {
	*event.Event

	// ChannelsToProcess defines channels to process
	ChannelsToProcess []ChannelToProcess
}
