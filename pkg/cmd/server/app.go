package server

import (
	"github.com/gin-gonic/gin"
)

type Options struct {
	Version string

	Env string

	Addr     string
	GinDebug bool
}

type App struct {
	Version string

	Addr string
	Env  string

	Router *gin.Engine
}

func SetupApp(options *Options) (*App, error) {
	if options == nil {
		options = &Options{}
	}

	app := &App{}
	app.Version = options.Version
	app.Addr = options.Addr
	app.Env = options.Env

	if options.GinDebug {
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

	r.GET("/", app.Index)
	r.Any("/healthz", app.Healthz)
	r.GET("/api/v1/healthz", app.Healthz)
	r.POST("/api/v1/callback-channel-test", app.CallbackChannelTest)

	r.HandleMethodNotAllowed = true

	app.Router = r

	return app, nil
}

func (app *App) Serve() error {
	return app.Router.Run(app.Addr)
}
