package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
	"net/http"
)

// ChannelCallback is the callback channel
type ChannelCallback struct {
	Client *http.Client
	URL    string
}

// Handle implements Channel
func (c *ChannelCallback) Handle(e *event.Event) error {
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

func NewChannelCallback(url string) (*ChannelCallback, error) {
	return &ChannelCallback{
		Client: http.DefaultClient,
		URL:    url,
	}, nil
}
