package channels

import (
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
}

func (c *ChannelTelegram) Handle(_ *event.Event) error {
	panic("wip")
}
