package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/controller"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/signals"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
	"time"
)

type controllerOptions struct {
	GlobalOptions     *genericoptions.GlobalOptions
	DatabaseOptions   *genericoptions.DatabaseOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions

	DefaultWorkers      int
	DefaultMaxRetries   int
	DefaultChannelNames []string
	MinResyncPeriod     time.Duration
}

func getControllerOptions() *controllerOptions {
	o := &controllerOptions{
		GlobalOptions:       genericoptions.GetGlobalOptions(),
		DatabaseOptions:     genericoptions.GetDatabaseOptions(),
		KubeconfigOptions:   genericoptions.GetKubeconfigOptions(),
		DefaultWorkers:      viper.GetInt("default-workers"),
		DefaultMaxRetries:   viper.GetInt("default-max-retries"),
		DefaultChannelNames: viper.GetStringSlice("default-channel-names"),
		MinResyncPeriod:     viper.GetDuration("min-resync-period"),
	}
	return o
}

func newControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Run controller, watch events and persistent into database (only one instance should be running)",
		Run: func(cmd *cobra.Command, args []string) {
			o := getControllerOptions()

			c, err := controller.Setup(controller.Config{
				DBDialect:           o.DatabaseOptions.DBDialect,
				DBArgs:              o.DatabaseOptions.DBArgs,
				Kubeconfig:          o.KubeconfigOptions.Kubeconfig,
				DefaultWorkers:      o.DefaultWorkers,
				DefaultMaxRetries:   o.DefaultMaxRetries,
				DefaultChannelNames: o.DefaultChannelNames,
				MinResyncPeriod:     o.MinResyncPeriod,
			})
			if err != nil {
				klog.Exit(errors.Wrap(err, "setup controller error"))
			}

			stopCh := signals.SetupSignalHandler()

			if err := c.Start(stopCh); err != nil {
				klog.Exit(errors.Wrap(err, "start controller error"))
			}
			defer c.Shutdown()

			<-stopCh
		},
	}

	f := cmd.PersistentFlags()
	f.Int("default-workers", 3, "default workers")
	f.Int("default-max-retries", 3, "default max retries")
	f.StringSlice("default-channel-names", nil, "default channel names")
	f.Duration("min-resync-period", 12*time.Hour, "min resync period (from n to 2n)")
	genericoptions.AddDatabaseFlags(f)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
