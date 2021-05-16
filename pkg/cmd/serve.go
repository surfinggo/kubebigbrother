package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/server"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
)

type ServeOptions struct {
	GlobalOptions   *genericoptions.GlobalOptions
	DatabaseOptions *genericoptions.DatabaseOptions

	Addr string
}

func GetServeOptions() *ServeOptions {
	o := &ServeOptions{
		GlobalOptions:   genericoptions.GetGlobalOptions(),
		DatabaseOptions: genericoptions.GetDatabaseOptions(),
		Addr:            viper.GetString("addr"),
	}
	return o
}

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the server to serve backend APIs",
		Run: func(cmd *cobra.Command, args []string) {
			o := GetServeOptions()
			app, err := server.SetupApp(&server.Options{
				Version: Version,
				Env:     o.GlobalOptions.Env,
				Addr:    o.Addr,
			})
			if err != nil {
				klog.Fatal(errors.Wrap(err, "setup app error"))
			}
			klog.Infof("env: %s", app.Env)
			klog.Infof("listening on: %s", app.Addr)
			err = app.Serve()
			if err != nil {
				klog.Fatal(err)
			}
		},
	}

	f := cmd.PersistentFlags()
	f.String("addr", "0.0.0.0:8984", "serving address")
	genericoptions.AddDatabaseFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
