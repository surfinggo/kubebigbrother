package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/cmd/genericoptions"
	"github.com/spongeprojects/kubebigbrother/pkg/gormdb"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
)

type queryOptions struct {
	GlobalOptions   *genericoptions.GlobalOptions
	DatabaseOptions *genericoptions.DatabaseOptions
}

func getQueryOptions() *queryOptions {
	o := &queryOptions{
		GlobalOptions:   genericoptions.GetGlobalOptions(),
		DatabaseOptions: genericoptions.GetDatabaseOptions(),
	}
	return o
}

func newQueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query event history",
		Run: func(cmd *cobra.Command, args []string) {
			o := getQueryOptions()

			db, err := gormdb.New(o.DatabaseOptions.DBDialect, o.DatabaseOptions.DBArgs)
			if err != nil {
				klog.Exit(errors.Wrap(err, "connect to db error"))
			}

			store := event_store.New(db)
			events, err := store.List(event_store.ListOptions{})
			if err != nil {
				klog.Exit(errors.Wrap(err, "list events error"))
			}
			if len(events) == 0 {
				fmt.Println("nothing")
			}
			for _, event := range events {
				fmt.Printf("ID: %d, %s\n", event.ID, event.Obj)
			}
		},
	}

	f := cmd.PersistentFlags()
	genericoptions.AddDatabaseFlags(f)
	magicconch.Must(viper.BindPFlags(f))

	return cmd
}
