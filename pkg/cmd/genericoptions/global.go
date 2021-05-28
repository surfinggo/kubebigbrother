package genericoptions

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/crumbs"
	"github.com/spongeprojects/magicconch"
	"k8s.io/klog/v2"
	"os"
)

// GlobalOptions is global options
type GlobalOptions struct {
	Env    string
	Config string
}

// IsDebugging checks whether env is debug
func (o GlobalOptions) IsDebugging() bool {
	return o.Env == crumbs.EnvDebug
}

// GetGlobalOptions get global options from viper flags
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

// AddGlobalFlags adds global flags to flag set
func AddGlobalFlags(f *pflag.FlagSet) {
	f.String("env", magicconch.Getenv("ENV", crumbs.EnvDebug), "environment")
	f.StringP("config", "c", crumbs.DefaultConfigFile, "path to config file (klog flags are not loaded from file, like -v)")
	addKlogFlags(f)
}
