package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Index
func (app *App) Index(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, current running version: %s", app.Version))
}

// Healthz, health check
func (app *App) Healthz(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, current running version: %s", app.Version))
}
