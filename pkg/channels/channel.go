package channels

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelName is name of channel
type ChannelName string

// ChannelType is type of channel
type ChannelType string

const (
	ChannelTypeGroup    = "group"    // a group of other channels
	ChannelTypeTelegram = "telegram" // send message to Telegram
	ChannelTypeCallback = "callback" // send message to callback url
	ChannelTypePrint    = "print"    // write message to writer
)

// ChannelMap maps from ChannelName to Channel
type ChannelMap map[ChannelName]Channel

// Channel is interface of a channel
type Channel interface {
	// NewProcessData builds a new copy of data to process,
	// for example, NewProcessData returns []recipient in Telegram channel，
	// if any error occurs and the processing is retried,
	// the channel can know which recipients have already been noticed successfully.
	NewProcessData() interface{}

	// Handle handles an event
	Handle(e *event.Event, processData interface{}) error
}
