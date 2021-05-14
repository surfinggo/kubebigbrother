package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/controller"
	"github.com/spongeprojects/kubebigbrother/pkg/fileorcreate"
	"github.com/spongeprojects/kubebigbrother/pkg/informers"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/magicconch"
	"os"
	"os/signal"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Run controller, watch events and persistent into database (only one instance should be running)",
	Run: func(cmd *cobra.Command, args []string) {
		informersConfigPath := viper.GetString("informers-config")
		err := fileorcreate.Ensure(informersConfigPath, InformersConfigFileTemplate)
		if err != nil {
			log.Error(errors.Wrap(err, "apply informers config template error"))
		}

		informersConfig, err := informers.LoadConfigFromFile(informersConfigPath)
		if err != nil {
			log.Fatal(errors.Wrap(err, "informers.LoadConfigFromFile error"))
		}
		controller, err := controller.Setup(controller.Options{
			KubeConfig:      viper.GetString("kubeconfig"),
			InformersConfig: informersConfig,
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "setup controller error"))
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

		controller.Start(stopCh)

		<-stopCh
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)

	f := controllerCmd.PersistentFlags()
	f.String("kubeconfig", defaultKubeconfig, "path to kubeconfig file")
	f.String("informers-config", DefaultInformersConfigFile, "path to informers config")

	magicconch.Must(viper.BindPFlags(f))
}
