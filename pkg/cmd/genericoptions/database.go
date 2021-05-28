package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DatabaseOptions is database options
type DatabaseOptions struct {
	DBDialect string
	DBArgs    string
}

// GetDatabaseOptions gets database options from viper flags
func GetDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		DBDialect: viper.GetString("db-dialect"),
		DBArgs:    viper.GetString("db-args"),
	}
}

// AddDatabaseFlags adds database flags to flag set
func AddDatabaseFlags(fs *pflag.FlagSet) {
	fs.String("db-dialect", "sqlite", "database dialect [mysql, postgres, sqlite]")
	fs.String("db-args", "", "database args")
}
