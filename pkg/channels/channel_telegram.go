package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	spg "github.com/spongeprojects/client-go/api/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
	Client      *http.Client
	Token       string
	ChatIDs     []string
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// NewEventProcessContext implements Channel
func (c *ChannelTelegram) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  c.ChatIDs,
	}
}

// TelegramMessage represents a Telegram message
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *ChannelTelegram) sendToRecipient(chatID string, text string) error {
	message := TelegramMessage{
		ChatID: chatID,
		Text:   text,
	}

	sendURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.Token)

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(message); err != nil {
		return errors.Wrap(err, "json encode error")
	}

	resp, err := c.Client.Post(sendURL, "application/json", body)
	if err != nil {
		return errors.Wrap(err, "send request error")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			klog.Warning(errors.Wrap(err, "close body error"))
		}
	}()
	if resp.StatusCode != 200 {
		return errors.Errorf("non-200 code returned: %d", resp.StatusCode)
	}
	return nil
}

// Handle implements Channel
func (c *ChannelTelegram) Handle(ctx *EventProcessContext) error {
	chatIDs := ctx.Data.([]string)

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
	message := buf.String()

	errs := make(map[string]error)
	for _, chatID := range chatIDs {
		if err := c.sendToRecipient(chatID, message); err != nil {
			errs[chatID] = err
		}
	}

	if len(errs) == 0 { // no error, no chatID left, everything works as expected
		ctx.Data = nil
		return nil
	}

	var chatIDsLeft []string
	var es []string
	for chatID, err := range errs {
		chatIDsLeft = append(chatIDsLeft, chatID)
		es = append(es, fmt.Sprintf("send to %s error: %s", chatID, err))
	}
	ctx.Data = chatIDsLeft
	return errors.Errorf("send Telegram message error: %s", strings.Join(es, ","))
}

// NewChannelTelegram creates new Telegram channel
func NewChannelTelegram(config *spg.ChannelTelegramConfig) (*ChannelTelegram, error) {
	if len(config.Token) < 40 {
		return nil, errors.New("invalid token, too short")
	}

	klog.V(2).Infof("Telegram token: %s...", config.Token[:15])

	var httpClient *http.Client
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", config.Proxy)
		}

		klog.V(2).Infof("connect to Telegram via proxy: %s", proxyUrl)

		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		httpClient = http.DefaultClient
	}

	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		config.AddedTemplate, config.DeletedTemplate, config.UpdatedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	var chatIDs []string
	for _, chatID := range config.ChatIDs {
		chatIDs = append(chatIDs, chatID)
	}

	return &ChannelTelegram{
		Client:      httpClient,
		Token:       config.Token,
		ChatIDs:     chatIDs,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
