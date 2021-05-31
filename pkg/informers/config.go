package informers

import (
	"github.com/pkg/errors"
)

type Config struct {
	KubeConfig string

	ConfigFile *ConfigFile
}

func (c *Config) Validate() error {
	if err := c.ConfigFile.Validate(); err != nil {
		return errors.Wrap(err, "invalid config")
	}
	return nil
}
