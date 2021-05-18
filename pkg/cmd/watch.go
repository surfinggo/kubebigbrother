package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
	"github.com/spongeprojects/kubebigbrother/pkg/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/signals"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
	"github.com/spongeprojects/kubebigbrother/staging/fileorcreate"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
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
				klog.Exit(errors.Wrap(err, "informers.LoadConfigFromFile error"))
			}

			w, err := watcher.Setup(watcher.Options{
				KubeConfig:      o.KubeconfigOptions.Kubeconfig,
				InformersConfig: informersConfig,
			})
			if err != nil {
				klog.Exit(errors.Wrap(err, "setup watcher error"))
			}

			stopCh := signals.SetupSignalHandler()

			if err := w.Start(stopCh); err != nil {
				klog.Exit(errors.Wrap(err, "start watcher error"))
			}
			defer w.Shutdown()

			<-stopCh
		},
	}

	f := cmd.PersistentFlags()
	genericoptions.AddInformersFlags(f)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
