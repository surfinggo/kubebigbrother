package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type InformersOptions struct {
	InformersConfig string
}

func GetInformersOptions() *InformersOptions {
	return &InformersOptions{
		InformersConfig: viper.GetString("informers-config"),
	}
}

func AddInformersFlags(fs *pflag.FlagSet, defaultInformersConfig string) {
	fs.String("informers-config", defaultInformersConfig, "path to informers config file")
}
