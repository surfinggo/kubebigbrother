package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
)

// InformersOptions is informers options
type InformersOptions struct {
	InformersConfig string
}

// GetInformersOptions get informers options from viper flags
func GetInformersOptions() *InformersOptions {
	return &InformersOptions{
		InformersConfig: viper.GetString("informers-config"),
	}
}

// AddInformersFlags adds informers flags to flag set
func AddInformersFlags(fs *pflag.FlagSet) {
	fs.String("informers-config", crumbs.DefaultInformersConfigFile, "path to informers config file")
}
