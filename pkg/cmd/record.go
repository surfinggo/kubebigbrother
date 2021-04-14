package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/watcher"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Run recorder, watch events and persistent into database",
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
	rootCmd.AddCommand(recordCmd)

	f := recordCmd.PersistentFlags()
	f.String("resource", "", "resource to watch")

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}
