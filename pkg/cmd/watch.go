package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
	"github.com/spongeprojects/magicconch"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Run watch to watch specific resource's event",
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
	rootCmd.AddCommand(watchCmd)

	f := watchCmd.PersistentFlags()
	f.String("resource", "", "resource to watch")
	f.String("kubeconfig", defaultKubeconfig, "kube config file path")

	magicconch.Must(viper.BindPFlags(f))
}
