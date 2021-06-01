package server

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Version string

	Addr                string
	Env                 string
	GinDebug            bool
	InformersConfigPath string
}

type App struct {
	Version string

	Addr                string
	Env                 string
	InformersConfigPath string

	Router *gin.Engine
}

func SetupApp(config *Config) (*App, error) {
	if config == nil {
		config = &Config{}
	}

	app := &App{}
	app.Version = config.Version
	app.Addr = config.Addr
	app.Env = config.Env
	app.InformersConfigPath = config.InformersConfigPath

	if config.GinDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)

	r.GET("/", app.HandlerIndex)
	r.Any("/healthz", app.HandlerHealthz)
	r.GET("/api/v1/healthz", app.HandlerHealthz)
	r.POST("/api/v1/callback-channel-test", app.HandlerCallbackChannelTest)
	r.GET("/api/v1/config", app.HandlerConfig)

	r.HandleMethodNotAllowed = true

	app.Router = r

	return app, nil
}

func (app *App) Serve() error {
	return app.Router.Run(app.Addr)
}
