package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"net/http"
	"text/template"
)

// ChannelCallbackConfig is config for ChannelCallback
type ChannelCallbackConfig struct {
	Method          string
	URL             string
	UseTemplate     bool
	AddedTemplate   string
	DeletedTemplate string
	UpdatedTemplate string
}

// ChannelCallback is the callback channel
type ChannelCallback struct {
	Client      *http.Client
	Method      string
	URL         string
	UseTemplate bool
	TmplAdded   *template.Template
	TmplDeleted *template.Template
	TmplUpdated *template.Template
}

// NewEventProcessContext implements Channel
func (c *ChannelCallback) NewEventProcessContext(e *event.Event) *EventProcessContext {
	return &EventProcessContext{
		Event: e,
		Data:  nil,
	}
}

// Handle implements Channel
func (c *ChannelCallback) Handle(ctx *EventProcessContext) error {
	body := &bytes.Buffer{}
	if c.UseTemplate {
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
		if err := t.Execute(body, ctx.Event); err != nil {
			return errors.Wrap(err, "execute template error")
		}
	} else {
		if err := json.NewEncoder(body).Encode(ctx.Event); err != nil {
			return errors.Wrap(err, "json encode error")
		}
	}
	req, err := http.NewRequest(c.Method, c.URL, body)
	if err != nil {
		return errors.Wrap(err, "build request error")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request error")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("non-200 code returned: %d", resp.StatusCode)
	}
	return nil
}

// NewChannelCallback creates callback channel
func NewChannelCallback(config *ChannelCallbackConfig) (*ChannelCallback, error) {
	if config.Method == "" {
		config.Method = "POST"
	}

	tmplAdded, tmplDeleted, tmplUpdated, err := parseTemplates(
		config.AddedTemplate, config.DeletedTemplate, config.UpdatedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse template error")
	}

	return &ChannelCallback{
		Client:      http.DefaultClient,
		Method:      config.Method,
		URL:         config.URL,
		UseTemplate: config.UseTemplate,
		TmplAdded:   tmplAdded,
		TmplDeleted: tmplDeleted,
		TmplUpdated: tmplUpdated,
	}, nil
}
