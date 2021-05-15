package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/fileorcreate"
	"github.com/spongeprojects/kubebigbrother/pkg/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
)

type WatchOptions struct {
	GlobalOptions     *genericoptions.GlobalOptions
	DatabaseOptions   *genericoptions.DatabaseOptions
	InformersOptions  *genericoptions.InformersOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions
}

func NewWatchOptions() *WatchOptions {
	o := &WatchOptions{
		GlobalOptions:     genericoptions.GetGlobalOptions(),
		DatabaseOptions:   genericoptions.GetDatabaseOptions(),
		InformersOptions:  genericoptions.GetInformersOptions(),
		KubeconfigOptions: genericoptions.GetKubeconfigOptions(),
	}
	return o
}

func NewWatchCommand() *cobra.Command {
	o := NewWatchOptions()

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Run watch to watch specific resource's event",
		Run: func(cmd *cobra.Command, args []string) {
			informersConfigPath := o.InformersOptions.InformersConfig
			err := fileorcreate.Ensure(informersConfigPath, InformersConfigFileTemplate)
			if err != nil {
				klog.Error(errors.Wrap(err, "apply informers config template error"))
			}

			informersConfig, err := informers.LoadConfigFromFile(informersConfigPath)
			if err != nil {
				klog.Fatal(errors.Wrap(err, "informers.LoadConfigFromFile error"))
			}
			w, err := watcher.Setup(watcher.Options{
				KubeConfig:      o.KubeconfigOptions.Kubeconfig,
				InformersConfig: informersConfig,
			})
			if err != nil {
				klog.Fatal(errors.Wrap(err, "setup watcher error"))
			}

			stopCh := make(chan struct{})

			// Ctrl+C
			interrupted := make(chan os.Signal)
			signal.Notify(interrupted, os.Interrupt)

			go func() {
				<-interrupted
				close(stopCh)
				<-interrupted // exit when interrupted again
				os.Exit(1)
			}()

			w.Start(stopCh)

			<-stopCh
		},
	}

	f := cmd.PersistentFlags()
	genericoptions.AddDatabaseFlags(f)
	genericoptions.AddInformersFlags(f, DefaultInformersConfigFile)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
