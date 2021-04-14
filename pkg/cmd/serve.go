package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := server.SetupApp(&server.Options{
			Version: Version,
			Env:     viper.GetString("env"),
			Addr:    viper.GetString("addr"),
		})
		if err != nil {
			return errors.Wrap(err, "setup app error")
		}
		return app.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	f := serveCmd.PersistentFlags()
	f.String("addr", "0.0.0.0:1949", "serving address")

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}
