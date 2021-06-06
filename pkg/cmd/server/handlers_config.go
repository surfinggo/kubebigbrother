package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

// HandlerConfig returns the currently used config
func (app *App) HandlerConfig(c *gin.Context) {
	channels, err := app.ChannelLister.List(labels.Everything())
	if err != nil {
		app.handle(c, errors.Wrap(err, "list channels error"))
	}

	watchers, err := app.WatcherLister.List(labels.Everything())
	if err != nil {
		app.handle(c, errors.Wrap(err, "list watchers error"))
	}

	clusterwatchers, err := app.ClusterWatcherLister.List(labels.Everything())
	if err != nil {
		app.handle(c, errors.Wrap(err, "list clusterwatchers error"))
	}

	c.JSON(200, gin.H{
		"channels":        channels,
		"watchers":        watchers,
		"clusterwatchers": clusterwatchers,
	})
}
