package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
)

type Config struct {
	KubeConfig string

	ConfigFile *ConfigFile

	SaveEvent  bool
	EventStore event_store.Interface
}

func (c *Config) Validate() error {
	if c.SaveEvent && c.EventStore == nil {
		return errors.New("event store cannot be nil when SaveEvent=true")
	}
	if err := c.ConfigFile.Validate(); err != nil {
		return errors.Wrap(err, "invalid config")
	}
	return nil
}
