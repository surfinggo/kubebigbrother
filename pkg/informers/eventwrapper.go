package informers

import (
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

type eventWrapper struct {
	*event.Event

	// ChannelsToProcess is channels to process, each key represents a channel,
	// keys are deleted from the map once the channel process successfully.
	//
	// The value of the map is data used by the channel,
	// for example, Telegram channel stores []recipient in the value,
	// recipients are deleted from the slice once message send successfully,
	// if any error occurs and the processing is retried,
	// it can know which recipients have already been noticed successfully.
	ChannelsToProcess map[channels.ChannelName]interface{}
}
