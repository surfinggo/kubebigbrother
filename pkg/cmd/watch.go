package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/watcher"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/style"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/signals"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
)

type watchOptions struct {
	GlobalOptions     *genericoptions.GlobalOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions
}

func getWatchOptions() *watchOptions {
	o := &watchOptions{
		GlobalOptions:     genericoptions.GetGlobalOptions(),
		KubeconfigOptions: genericoptions.GetKubeconfigOptions(),
	}
	return o
}

func newWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch events lively",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(style.Fg(style.Warning, ""+
				"--------------------------------------------------\n"+
				"|  Watch should only be used for debugging.      |\n"+
				"|  In watch mode, all channels will be replaced  |\n"+
				"|  by a single \"print to stdout\" channel.        |\n"+
				"--------------------------------------------------"))

			o := getWatchOptions()

			w, err := watcher.Setup(watcher.Config{
				Kubeconfig: o.KubeconfigOptions.Kubeconfig,
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
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
