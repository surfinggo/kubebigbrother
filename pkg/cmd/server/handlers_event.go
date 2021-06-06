package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"github.com/spongeprojects/magicconch"
	"gopkg.in/yaml.v3"
)

// HandlerEventList queries events
func (app *App) HandlerEventList(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")
	var informerName string
	if namespace != "" {
		informerName = models.WatcherInformerName(namespace, name)
	} else if name != "" {
		informerName = models.ClusterWatcherInformerName(name)
	}
	events, err := app.EventStore.List(event_store.ListOptions{
		InformerName: informerName,
		Q:            c.Query("q"),
		After:        magicconch.StringToUint(c.Query("after")),
	})
	if err != nil {
		app.handle(c, errors.Wrap(err, "list events error"))
		return
	}

	for i := range events {
		evt := events[i]
		evt.Obj = nil
		evt.OldObj = nil
		events[i] = evt
	}

	c.JSON(200, gin.H{
		"events": events,
	})
	return
}

// HandlerEvent get event by id
func (app *App) HandlerEvent(c *gin.Context) {
	event, err := app.EventStore.Find(magicconch.StringToUint(c.Param("id")))
	if err != nil {
		app.handle(c, errors.Wrap(err, "find events error"))
		return
	}

	obj := event.GetObj()
	oldObj := event.GetOldObj()

	var objYamlBytes, oldObjYamlBytes []byte

	if obj != nil {
		objYamlBytes, err = yaml.Marshal(obj)
		if err != nil {
			app.handle(c, errors.Wrap(err, "yaml marshal error"))
		}
	}
	if oldObj != nil {
		oldObjYamlBytes, err = yaml.Marshal(oldObj)
		if err != nil {
			app.handle(c, errors.Wrap(err, "yaml marshal error"))
		}
	}

	event.Obj = nil
	event.OldObj = nil

	eventYamlBytes, err := yaml.Marshal(event)
	if err != nil {
		app.handle(c, errors.Wrap(err, "yaml marshal error"))
	}

	c.JSON(200, gin.H{
		"event":        event,
		"event_yaml":   eventYamlBytes,
		"obj":          obj,
		"old_obj":      oldObj,
		"obj_yaml":     objYamlBytes,
		"old_obj_yaml": oldObjYamlBytes,
	})
	return
}
