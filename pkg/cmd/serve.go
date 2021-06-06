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
	GlobalOptions     *genericoptions.GlobalOptions
	DatabaseOptions   *genericoptions.DatabaseOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions

	GinDebug bool
	Addr     string
}

func getServeOptions() *serveOptions {
	o := &serveOptions{
		GlobalOptions:     genericoptions.GetGlobalOptions(),
		DatabaseOptions:   genericoptions.GetDatabaseOptions(),
		KubeconfigOptions: genericoptions.GetKubeconfigOptions(),
		GinDebug:          viper.GetBool("gin-debug"),
		Addr:              viper.GetString("addr"),
	}
	return o
}

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the server to serve backend APIs",
		Run: func(cmd *cobra.Command, args []string) {
			o := getServeOptions()

			app, err := server.SetupApp(&server.Config{
				Env:        o.GlobalOptions.Env,
				Version:    Version,
				Addr:       o.Addr,
				DBDialect:  o.DatabaseOptions.DBDialect,
				DBArgs:     o.DatabaseOptions.DBArgs,
				GinDebug:   o.GinDebug,
				Kubeconfig: o.KubeconfigOptions.Kubeconfig,
			})
			if err != nil {
				klog.Exit(errors.Wrap(err, "setup app error"))
			}

			klog.Infof("env: %s", app.Env)
			klog.Infof("listening on: %s", app.Addr)

			stopCh := make(chan struct{})
			defer close(stopCh)

			err = app.Run(stopCh)
			if err != nil {
				klog.Exit(err)
			}
		},
	}

	f := cmd.PersistentFlags()
	f.String("addr", "0.0.0.0:8984", "serving address")
	f.Bool("gin-debug", false, "enable gin debug mode")
	genericoptions.AddDatabaseFlags(f)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
