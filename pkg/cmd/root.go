package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/fileorcreate"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/homedir"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/magicconch"
	"math/rand"
	"os"
	"path"
	"time"
)

var Version = "unknown"

const (
	EnvDebug = "debug"

	DefaultConfigFile = "config/config.local.yaml"

	ConfigFileTemplate = "config/config.tmpl.yaml"

	DefaultInformersConfigFile = "config/informers-config.local.yaml"

	InformersConfigFileTemplate = "config/informers-config.tmpl.yaml"
)

var (
	cfgFile string
	env     string

	defaultKubeconfig = magicconch.Getenv("KUBECONFIG", path.Join(homedir.HomeDir(), ".kube", "config"))
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kbb",
	Short: "`kbb` command line tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	cobra.OnInitialize(initConfig)

	f := rootCmd.PersistentFlags()
	f.StringVarP(&cfgFile, "config", "c", DefaultConfigFile, "config file")
	f.StringVarP(&env, "env", "e", os.Getenv("ENV"), "environment")

	magicconch.Must(viper.BindPFlags(f))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("KBB")
	viper.AutomaticEnv()

	if env == "" {
		env = EnvDebug
	}

	if env == EnvDebug {
		if _, exist := os.LookupEnv("LOG_LEVEL"); !exist {
			log.Logger.SetLevel(logrus.DebugLevel)
		}
	}

	err := fileorcreate.Ensure(cfgFile, ConfigFileTemplate)
	if err != nil {
		log.Error(errors.Wrap(err, "apply config template error"))
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		err := viper.ReadInConfig()
		if err != nil {
			log.Warn(errors.Wrapf(err, "read in config error, file: %s", viper.ConfigFileUsed()))
		} else {
			log.Infof("using config file: %s", viper.ConfigFileUsed())
		}
	} else {
		log.Info("config file not specified, not reading from file")
	}
}
