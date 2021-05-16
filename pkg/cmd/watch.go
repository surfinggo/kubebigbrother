package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
	"github.com/spongeprojects/kubebigbrother/pkg/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
	"github.com/spongeprojects/kubebigbrother/staging/fileorcreate"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
)

type WatchOptions struct {
	GlobalOptions     *genericoptions.GlobalOptions
	InformersOptions  *genericoptions.InformersOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions
}

func GetWatchOptions() *WatchOptions {
	o := &WatchOptions{
		GlobalOptions:     genericoptions.GetGlobalOptions(),
		InformersOptions:  genericoptions.GetInformersOptions(),
		KubeconfigOptions: genericoptions.GetKubeconfigOptions(),
	}
	return o
}

func NewWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Run watch to watch specific resource's event",
		Run: func(cmd *cobra.Command, args []string) {
			o := GetWatchOptions()

			informersConfigPath := o.InformersOptions.InformersConfig

			if o.GlobalOptions.IsDebugging() {
				err := fileorcreate.Ensure(informersConfigPath, crumbs.InformersConfigFileTemplate)
				if err != nil {
					klog.Error(errors.Wrap(err, "apply informers config template error"))
				}
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

			if err := w.Start(stopCh); err != nil {
				klog.Fatal(errors.Wrap(err, "start watcher error"))
			}

			<-stopCh
		},
	}

	f := cmd.PersistentFlags()
	genericoptions.AddInformersFlags(f)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
