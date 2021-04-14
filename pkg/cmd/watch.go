package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Run watcher",
	RunE: func(cmd *cobra.Command, args []string) error {
		watcher, err := watcher.Setup(watcher.Options{
			Env:        env,
			KubeConfig: viper.GetString("kube-config"),
			Resource:   viper.GetString("resource"),
		})
		if err != nil {
			return errors.Wrap(err, "setup watcher error")
		}
		return watcher.Start()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	f := watchCmd.PersistentFlags()
	f.String("resource", "", "resource to watch")

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}
