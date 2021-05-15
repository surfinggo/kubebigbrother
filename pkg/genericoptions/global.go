package genericoptions

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
)

type GlobalOptions struct {
	Env    string
	Config string
}

func GetGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		Env:    viper.GetString("env"),
		Config: viper.GetString("config"),
	}
}

// addKlogFlags adds flags from k8s.io/klog
func addKlogFlags(fs *pflag.FlagSet) {
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	normalizeFunc := fs.GetNormalizeFunc()
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = string(normalizeFunc(fs, fl.Name))
		fs.AddGoFlag(fl)
	})
}

func AddGlobalFlags(f *pflag.FlagSet, defaultEnv, defaultConfigFile string) {
	f.String("env", defaultEnv, "environment")
	f.StringP("config", "c", defaultConfigFile, "path to config file (klog flags are not loaded from file, like -v)")
	addKlogFlags(f)
}
