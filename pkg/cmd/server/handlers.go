package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// Index the index endpoint
func (app *App) Index(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, current running version: %s", app.Version))
}

// Healthz is the health check endpoint
func (app *App) Healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": app.Version,
	})
}

// CallbackChannelTest is the endpoint for callback channel test
func (app *App) CallbackChannelTest(c *gin.Context) {
	_, _ = os.Stdout.Write([]byte("----------callback channel test, data received:----------\n"))
	_, _ = io.Copy(os.Stdout, c.Request.Body)
	_, _ = os.Stdout.Write([]byte("\n----------callback channel test, end data----------\n"))
	c.String(200, "OK")
}
