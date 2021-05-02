package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
	"github.com/spongeprojects/magicconch"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Run recorder, watch events and persistent into database (only one instance should be running)",
	Run: func(cmd *cobra.Command, args []string) {
		watcher, err := watcher.Setup(watcher.Options{
			Env:        env,
			KubeConfig: viper.GetString("kubeconfig"),
			Resource:   viper.GetString("resource"),
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "setup watcher error"))
		}
		err = watcher.Start()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)

	f := recordCmd.PersistentFlags()
	f.String("resource", "", "resource to watch")
	f.String("kubeconfig", defaultKubeconfig, "kube config file path")

	magicconch.Must(viper.BindPFlags(f))
}
