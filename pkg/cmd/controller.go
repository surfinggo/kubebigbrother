package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/controller"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
	"github.com/spongeprojects/kubebigbrother/pkg/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/signals"
	"github.com/spongeprojects/kubebigbrother/staging/fileorcreate"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
)

type ControllerOptions struct {
	GlobalOptions     *genericoptions.GlobalOptions
	DatabaseOptions   *genericoptions.DatabaseOptions
	InformersOptions  *genericoptions.InformersOptions
	KubeconfigOptions *genericoptions.KubeconfigOptions
}

func GetControllerOptions() *ControllerOptions {
	o := &ControllerOptions{
		GlobalOptions:     genericoptions.GetGlobalOptions(),
		DatabaseOptions:   genericoptions.GetDatabaseOptions(),
		InformersOptions:  genericoptions.GetInformersOptions(),
		KubeconfigOptions: genericoptions.GetKubeconfigOptions(),
	}
	return o
}

func NewControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Run controller, watch events and persistent into database (only one instance should be running)",
		Run: func(cmd *cobra.Command, args []string) {
			o := GetControllerOptions()

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
			c, err := controller.Setup(controller.Options{
				DBDialect:       o.DatabaseOptions.DBDialect,
				DBArgs:          o.DatabaseOptions.DBDialect,
				KubeConfig:      o.KubeconfigOptions.Kubeconfig,
				InformersConfig: informersConfig,
			})
			if err != nil {
				klog.Fatal(errors.Wrap(err, "setup controller error"))
			}

			stopCh := signals.SetupSignalHandler()

			if err := c.Start(stopCh); err != nil {
				klog.Fatal(errors.Wrap(err, "start controller error"))
			}

			<-stopCh
		},
	}

	f := cmd.PersistentFlags()
	genericoptions.AddDatabaseFlags(f)
	genericoptions.AddInformersFlags(f)
	genericoptions.AddKubeconfigFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
