package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/server"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
)

type serveOptions struct {
	GlobalOptions    *genericoptions.GlobalOptions
	DatabaseOptions  *genericoptions.DatabaseOptions
	InformersOptions *genericoptions.InformersOptions

	Addr string
}

func getServeOptions() *serveOptions {
	o := &serveOptions{
		GlobalOptions:    genericoptions.GetGlobalOptions(),
		DatabaseOptions:  genericoptions.GetDatabaseOptions(),
		InformersOptions: genericoptions.GetInformersOptions(),
		Addr:             viper.GetString("addr"),
	}
	return o
}

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the server to serve backend APIs",
		Run: func(cmd *cobra.Command, args []string) {
			o := getServeOptions()

			informersConfigPath := o.InformersOptions.InformersConfig

			app, err := server.SetupApp(&server.Config{
				Version:             Version,
				Env:                 o.GlobalOptions.Env,
				Addr:                o.Addr,
				InformersConfigPath: informersConfigPath,
			})
			if err != nil {
				klog.Exit(errors.Wrap(err, "setup app error"))
			}

			klog.Infof("env: %s", app.Env)
			klog.Infof("listening on: %s", app.Addr)

			err = app.Serve()
			if err != nil {
				klog.Exit(err)
			}
		},
	}

	f := cmd.PersistentFlags()
	f.String("addr", "0.0.0.0:8984", "serving address")
	genericoptions.AddDatabaseFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
