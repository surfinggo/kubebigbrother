package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/kubebigbrother/pkg/server"
	"github.com/spongeprojects/magicconch"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server to serve backend APIs",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := server.SetupApp(&server.Options{
			Version: Version,
			Env:     env,
			Addr:    viper.GetString("addr"),
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "setup app error"))
		}
		err = app.Serve()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	f := serveCmd.PersistentFlags()
	f.String("addr", "0.0.0.0:8984", "serving address")

	magicconch.Must(viper.BindPFlags(f))
}
