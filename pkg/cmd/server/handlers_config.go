package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"gopkg.in/yaml.v3"
)

// HandlerConfig returns the currently used config
func (app *App) HandlerConfig(c *gin.Context) {
	informersConfig, err := informers.LoadConfigFromFile(app.InformersConfigPath)
	if err != nil {
		app.handle(c, errors.Wrap(err, "load config file error"))
		return
	}

	yamlBytes, err := yaml.Marshal(informersConfig)
	if err != nil {
		app.handle(c, errors.Wrap(err, "yaml marshal error"))
		return
	}

	c.JSON(200, gin.H{
		"json": informersConfig,
		"yaml": yamlBytes,
	})
}
