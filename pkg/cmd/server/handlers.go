package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// HandlerIndex the index endpoint
func (app *App) HandlerIndex(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, current running version: %s", app.Version))
}

// HandlerHealthz is the health check endpoint
func (app *App) HandlerHealthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": app.Version,
	})
}

// HandlerCallbackChannelTest is the endpoint for callback channel test
func (app *App) HandlerCallbackChannelTest(c *gin.Context) {
	_, _ = os.Stdout.Write([]byte("----------callback channel test, data received:----------\n"))
	_, _ = io.Copy(os.Stdout, c.Request.Body)
	_, _ = os.Stdout.Write([]byte("\n----------callback channel test, end data----------\n"))
	c.String(200, "OK")
}
