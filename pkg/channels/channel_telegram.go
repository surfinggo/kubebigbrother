package channels

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"github.com/spongeprojects/kubebigbrother/pkg/services/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
	"html/template"
	"strings"
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

// NewEventProcessContext implements Channel
func (c *ChannelTelegram) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  c.Recipients,
	}
}

// Handle implements Channel
func (c *ChannelTelegram) Handle(ctx *EventProcessContext) error {
	recipients := ctx.Data.([]tb.Recipient)

	var t *template.Template
	switch ctx.Event.Type {
	case event.TypeAdded:
		t = c.TmplAdded
	case event.TypeDeleted:
		t = c.TmplDeleted
	case event.TypeUpdated:
		t = c.TmplUpdated
	default:
		return errors.Errorf("unknown event type: %s", ctx.Event.Type)
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ctx.Event); err != nil {
		return errors.Wrap(err, "execute template error")
	}

	errs := make(map[tb.Recipient]error)
	for _, recipient := range recipients {
		_, err := c.Bot.Send(recipient, buf.String())
		if err != nil {
			errs[recipient] = err
		}
	}

	if len(errs) == 0 { // no error, no recipient left, everything works as expected
		ctx.Data = nil
		return nil
	}

	var recipientsLeft []tb.Recipient
	var es []string
	for recipient, err := range errs {
		recipientsLeft = append(recipientsLeft, recipient)
		es = append(es, fmt.Sprintf("send to %s error: %s", recipient, err))
	}
	ctx.Data = recipientsLeft
	return errors.Errorf("send Telegram message error: %s", strings.Join(es, ","))
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
		return nil, errors.Wrap(err, "create Telegram client error")
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
