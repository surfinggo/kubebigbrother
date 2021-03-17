package watcher

import "github.com/pkg/errors"

type Options struct {
	Version string

	Env string

	Resource string
	GinDebug bool
}

type App struct {
	Version string

	Addr string

	Controller interface{}
}

func SetupWatcher(options *Options) (*App, error) {
	if options == nil {
		options = &Options{}
	}

	app := &App{}
	app.Addr = options.Resource
	app.Version = options.Version
	app.Controller = nil

	return app, nil
}

func (app *App) Start() error {
	return errors.Wrap(err, "not implemented")
}
