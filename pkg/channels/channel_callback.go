package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"net/http"
)

// ChannelCallbackConfig is config for ChannelCallback
type ChannelCallbackConfig struct {
	URL string
}

// ChannelCallback is the callback channel
type ChannelCallback struct {
	Client *http.Client
	URL    string
}

// NewProcessData implements Channel
func (c *ChannelCallback) NewProcessData() interface{} {
	return nil
}

// Handle implements Channel
func (c *ChannelCallback) Handle(e *event.Event, _ interface{}) error {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(e)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	resp, err := c.Client.Post(c.URL, "application/json", body)
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
	return &ChannelCallback{
		Client: http.DefaultClient,
		URL:    config.URL,
	}, nil
}
