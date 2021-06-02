package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HandlerEvent queries events
func (app *App) HandlerEvent(c *gin.Context) {
	if c.Query("informerName") != "" {
		events, err := app.EventStore.ListByInformer(c.Query("informerName"))
		if err != nil {
			app.handle(c, errors.Wrap(err, "list events error"))
			return
		}

		c.JSON(200, gin.H{
			"events": events,
		})
		return
	}
	events, err := app.EventStore.Search(c.Query("q"))
	if err != nil {
		app.handle(c, errors.Wrap(err, "search events error"))
		return
	}

	c.JSON(200, gin.H{
		"events": events,
	})
	return
}
