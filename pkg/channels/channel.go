package channels

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelType is type of channel
type ChannelType string

const (
	ChannelTypeCallback = "callback" // send message to callback url
	ChannelTypeDingtalk = "dingtalk" // send message to Dingtalk url
	ChannelTypeFlock    = "flock"    // send message to Flock
	ChannelTypePrint    = "print"    // write message to writer
	ChannelTypeSlack    = "slack"    // send message to Slack
	ChannelTypeTelegram = "telegram" // send message to Telegram
)

// ChannelMap maps from string to Channel
type ChannelMap map[string]Channel

// Channel is interface of a channel
type Channel interface {
	// NewEventProcessContext builds a new copy of data to process for an event,
	// for example, NewEventProcessContext returns []chatID in Telegram channel，
	// if any error occurs and the processing is retried,
	// the channel can know which chatIDs have already been noticed successfully.
	NewEventProcessContext(e *event.Event) *EventProcessContext

	// Handle handles an event
	Handle(ctx *EventProcessContext) error
}

// EventProcessContext is the context of an event processing within a channel
type EventProcessContext struct {
	Event *event.Event

	// Data is the context data used by a channel during processing an event,
	// for example, Telegram channel stores []chatID in Data,
	// chatID are deleted from the slice once message send successfully,
	// if any error occurs and the processing is retried,
	// it can know which chatIDs have already been noticed successfully.
	Data interface{}
}
