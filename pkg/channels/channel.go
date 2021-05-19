package channels

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelName is name of channel
type ChannelName string

// ChannelType is type of channel
type ChannelType string

const (
	ChannelTypeGroup    = "group"
	ChannelTypeTelegram = "telegram"
	ChannelTypeCallback = "callback"
	ChannelTypePrint    = "print"
)

type ChannelMap map[ChannelName]Channel

type Channel interface {
	Handle(e *event.Event) error
}
