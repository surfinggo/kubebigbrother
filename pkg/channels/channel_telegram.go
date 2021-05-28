package channels

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/services/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
	"k8s.io/klog/v2"
)

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
	Bot        *tb.Bot
	Recipients []tb.Recipient
}

// Handle implements Channel
func (c *ChannelTelegram) Handle(e *event.Event) error {
	for _, recipient := range c.Recipients {
		_, err := c.Bot.Send(recipient, fmt.Sprintf("%s", e))
		if err != nil {
			klog.Warningf("send Telegram message error: %s", err)
		}
	}
	return nil
}

func NewChannelTelegram(token string, recipients []int64) (*ChannelTelegram, error) {
	bot, err := telegram.NewBot(token)
	if err != nil {
		return nil, errors.Wrap(err, "create bot error")
	}
	var tbRecipients []tb.Recipient
	for _, chatID := range recipients {
		tbRecipients = append(tbRecipients, &tb.Chat{ID: chatID})
	}
	return &ChannelTelegram{
		Bot:        bot,
		Recipients: tbRecipients,
	}, nil
}
