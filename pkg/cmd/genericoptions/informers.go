package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
)

type InformersOptions struct {
	InformersConfig string
}

func GetInformersOptions() *InformersOptions {
	return &InformersOptions{
		InformersConfig: viper.GetString("informers-config"),
	}
}

func AddInformersFlags(fs *pflag.FlagSet) {
	fs.String("informers-config", crumbs.DefaultInformersConfigFile, "path to informers config file")
}
