package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DatabaseOptions struct {
	DBDialect string
	DBArgs    string
}

func GetDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		DBDialect: viper.GetString("db-dialect"),
		DBArgs:    viper.GetString("db-args"),
	}
}

func AddDatabaseFlags(fs *pflag.FlagSet) {
	fs.String("db-dialect", "sqlite", "database dialect [mysql, postgres, sqlite]")
	fs.String("db-args", "", "database args")
}
