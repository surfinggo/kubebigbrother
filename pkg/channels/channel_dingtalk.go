package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	spg "github.com/spongeprojects/client-go/api/spongeprojects.com/v1alpha1"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"text/template"
)

// ChannelDingtalk is the callback channel
type ChannelDingtalk struct {
	Client      *http.Client
	WebhookURL  string
	AtMobiles   []string
	AtAll       bool
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// DingtalkMessageAt represents a Dingtalk message at (@) info
type DingtalkMessageAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

// DingtalkMessageText represents a Dingtalk message text
type DingtalkMessageText struct {
	Content string `json:"content"`
}

// DingtalkMessage represents a Dingtalk message
// ref: https://developers.dingtalk.com/document/app/custom-robot-access
type DingtalkMessage struct {
	At      DingtalkMessageAt   `json:"at"`
	Text    DingtalkMessageText `json:"text"`
	Msgtype string              `json:"msgtype"`
}

// NewEventProcessContext implements Channel
func (c *ChannelDingtalk) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  nil,
	}
}

// Handle implements Channel
func (c *ChannelDingtalk) Handle(ctx *EventProcessContext) error {
	buf := &bytes.Buffer{}
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

	if err := t.Execute(buf, ctx.Event); err != nil {
		return errors.Wrap(err, "execute template error")
	}

	message := DingtalkMessage{
		At: DingtalkMessageAt{
			AtMobiles: c.AtMobiles,
			IsAtAll:   c.AtAll,
		},
		Text: DingtalkMessageText{
			Content: buf.String(),
		},
		Msgtype: "text",
	}

	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(message); err != nil {
		return errors.Wrap(err, "json encode error")
	}

	resp, err := c.Client.Post(c.WebhookURL, "application/json", body)
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

// NewChannelDingtalk creates callback channel
func NewChannelDingtalk(config *spg.ChannelDingtalkConfig) (*ChannelDingtalk, error) {
	if len(config.WebhookURL) < 70 {
		return nil, errors.New("invalid url, too short")
	}

	klog.V(2).Infof("Dingtalk url: %s...", config.WebhookURL[:45])

	var httpClient *http.Client
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", config.Proxy)
		}

		klog.V(2).Infof("connect to Dingtalk via proxy: %s", proxyUrl)

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

	return &ChannelDingtalk{
		Client:      httpClient,
		WebhookURL:  config.WebhookURL,
		AtMobiles:   config.AtMobiles,
		AtAll:       config.AtAll,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
