package channels

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"text/template"
)

// ChannelSlackConfig is config for ChannelSlack
type ChannelSlackConfig struct {
	Token           string
	Proxy           string
	WebhookURL      string
	TitleTemplate   string
	AddedTemplate   string
	DeletedTemplate string
	UpdatedTemplate string
}

// ChannelSlack is the callback channel
type ChannelSlack struct {
	Client      *slack.Client // TODO: add Slack app support (not only webhooks)
	WebhookURL  string
	TmplTitle   *template.Template
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// NewEventProcessContext implements Channel
func (c *ChannelSlack) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  nil,
	}
}

// Handle implements Channel
// ref: https://api.slack.com/apps
func (c *ChannelSlack) Handle(ctx *EventProcessContext) error {
	titleBuf := &bytes.Buffer{}
	if err := c.TmplTitle.Execute(titleBuf, ctx.Event); err != nil {
		return errors.Wrap(err, "execute title template error")
	}
	title := titleBuf.String()

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

	err := slack.PostWebhook(c.WebhookURL, &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			{
				Color: ctx.Event.Color(),
				Title: title,
				Text:  buf.String(),
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "post message error")
	}

	return nil
}

// NewChannelSlack creates callback channel
func NewChannelSlack(config *ChannelSlackConfig) (*ChannelSlack, error) {
	var httpClient *http.Client
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", config.Proxy)
		}

		klog.V(2).Infof("connect to Slack via proxy: %s", proxyUrl)

		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		httpClient = http.DefaultClient
	}

	client := slack.New(config.Token, slack.OptionHTTPClient(httpClient))

	if config.TitleTemplate == "" {
		config.TitleTemplate = "New Event:"
		// context: event.Event
		//config.TitleTemplate = "New Event [{{.Type}}]:"
	}
	tmplTitle, err := template.New("").Parse(config.TitleTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse title template error")
	}

	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		config.AddedTemplate, config.DeletedTemplate, config.UpdatedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	return &ChannelSlack{
		Client:      client,
		WebhookURL:  config.WebhookURL,
		TmplTitle:   tmplTitle,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
