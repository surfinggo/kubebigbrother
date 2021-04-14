package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Query event history",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := gormdb.New(dbDialect, dbArgs)
		if err != nil {
			return errors.Wrap(err, "create db error")
		}
		store := event_store.New(db)
		events, err := store.List()
		if err != nil {
			return errors.Wrap(err, "list events error")
		}
		if len(events) == 0 {
			fmt.Println("nothing")
		}
		for _, event := range events {
			fmt.Printf("ID: %d, %s\n", event.ID, event.Description)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	f := historyCmd.PersistentFlags()
	f.String("resource", "", "resource to query")
	//f.String("kube-config", magicconch.Getenv("KUBECONFIG", os.Getenv("HOME")+"/.kube/config"), "kube config file path")

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}
