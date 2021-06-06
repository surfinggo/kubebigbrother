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

// ChannelFlock is the callback channel
type ChannelFlock struct {
	Client      *http.Client
	URL         string
	TmplTitle   *template.Template
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// FlockMessage represents a flock message
// ref: https://docs.flock.com/display/flockos/Message
type FlockMessage struct {
	Notification string                   `json:"notification"`
	Text         string                   `json:"text"`
	Attachments  []FlockMessageAttachment `json:"attachments"`
}

// FlockMessageAttachment represents a flock message attachment
// ref: https://docs.flock.com/display/flockos/Attachment
type FlockMessageAttachment struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

// NewEventProcessContext implements Channel
func (c *ChannelFlock) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  nil,
	}
}

// Handle implements Channel
func (c *ChannelFlock) Handle(ctx *EventProcessContext) error {
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

	message := FlockMessage{
		Text: title,
		Attachments: []FlockMessageAttachment{
			{
				Title: buf.String(),
				Color: ctx.Event.Color(),
			},
		},
	}

	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(message); err != nil {
		return errors.Wrap(err, "json encode error")
	}

	resp, err := c.Client.Post(c.URL, "application/json", body)
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

// NewChannelFlock creates callback channel
func NewChannelFlock(config *spg.ChannelFlockConfig) (*ChannelFlock, error) {
	if len(config.URL) < 70 {
		return nil, errors.New("invalid url, too short")
	}

	klog.V(2).Infof("Flock url: %s...", config.URL[:45])

	var httpClient *http.Client
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", config.Proxy)
		}

		klog.V(2).Infof("connect to Flock via proxy: %s", proxyUrl)

		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		httpClient = http.DefaultClient
	}

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

	return &ChannelFlock{
		Client:      httpClient,
		URL:         config.URL,
		TmplTitle:   tmplTitle,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
