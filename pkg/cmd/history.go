package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"github.com/spongeprojects/magicconch"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Query event history",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gormdb.New(viper.GetString("db-dialect"), viper.GetString("db-args"))
		if err != nil {
			log.Fatal(errors.Wrap(err, "connect to db error"))
		}
		store := event_store.New(db)
		events, err := store.List()
		if err != nil {
			log.Fatal(errors.Wrap(err, "list events error"))
		}
		if len(events) == 0 {
			fmt.Println("nothing")
		}
		for _, event := range events {
			fmt.Printf("ID: %d, %s\n", event.ID, event.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	f := historyCmd.PersistentFlags()
	f.String("resource", "", "resource to query")
	f.String("db-dialect", "sqlite", "database dialect [mysql, postgres, sqlite]")
	f.String("db-args", "", "database args")

	magicconch.Must(viper.BindPFlags(f))
}
