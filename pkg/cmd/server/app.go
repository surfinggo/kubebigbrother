package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
)

type Config struct {
	Version string

	Addr                string
	Env                 string
	DBDialect           string
	DBArgs              string
	GinDebug            bool
	InformersConfigPath string
}

type App struct {
	Version string

	Addr                string
	Env                 string
	InformersConfigPath string

	EventStore event_store.Interface
	Router     *gin.Engine
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

	db, err := gormdb.New(config.DBDialect, config.DBArgs)
	if err != nil {
		return nil, errors.Wrap(err, "create db instance error")
	}

	app.EventStore = event_store.New(db)

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
	r.GET("/api/v1/events", app.HandlerEventList)
	r.GET("/api/v1/events/:id", app.HandlerEvent)

	r.HandleMethodNotAllowed = true

	app.Router = r

	return app, nil
}

func (app *App) Serve() error {
	return app.Router.Run(app.Addr)
}
