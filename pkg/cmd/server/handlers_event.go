package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"github.com/spongeprojects/magicconch"
)

// HandlerEvent queries events
func (app *App) HandlerEvent(c *gin.Context) {
	events, err := app.EventStore.List(event_store.ListOptions{
		InformerName: c.Query("informerName"),
		Q:            c.Query("q"),
		After:        magicconch.StringToUint(c.Query("after")),
	})
	if err != nil {
		app.handle(c, errors.Wrap(err, "list events error"))
		return
	}

	c.JSON(200, gin.H{
		"events": events,
	})
	return
}
