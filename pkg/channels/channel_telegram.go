package channels

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/services/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
	"html/template"
	"k8s.io/klog/v2"
)

// ChannelTelegramConfig is config for ChannelTelegram
type ChannelTelegramConfig struct {
	Token           string
	Recipients      []int64
	Proxy           string
	AddedTemplate   string
	DeletedTemplate string
	UpdatedTemplate string
}

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
	Bot         *tb.Bot
	Recipients  []tb.Recipient
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// Handle implements Channel
func (c *ChannelTelegram) Handle(e *event.Event) error {
	var t *template.Template
	switch e.Type {
	case event.TypeAdded:
		t = c.TmplAdded
	case event.TypeDeleted:
		t = c.TmplDeleted
	case event.TypeUpdated:
		t = c.TmplUpdated
	default:
		return errors.Errorf("unknown event type: %s", e.Type)
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, e); err != nil {
		return errors.Wrap(err, "execute template error")
	}

	for _, recipient := range c.Recipients {
		_, err := c.Bot.Send(recipient, buf.String())
		if err != nil {
			// TODO: retry on failed recipients
			klog.Warningf("send Telegram message error: %s", err)
		}
	}
	return nil
}

// NewChannelTelegram creates new Telegram channel
func NewChannelTelegram(config *ChannelTelegramConfig) (*ChannelTelegram, error) {
	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		config.AddedTemplate, config.DeletedTemplate, config.UpdatedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	bot, err := telegram.NewBot(config.Token, config.Proxy)
	if err != nil {
		return nil, errors.Wrap(err, "create Telegram bot error")
	}

	var tbRecipients []tb.Recipient
	for _, chatID := range config.Recipients {
		tbRecipients = append(tbRecipients, &tb.Chat{ID: chatID})
	}

	return &ChannelTelegram{
		Bot:         bot,
		Recipients:  tbRecipients,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
