package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"github.com/spongeprojects/magicconch"
	"io"
	"math/rand"
	"os"
	"time"
)

var Version = "unknown"

var cfgFile string

const (
	EnvDebug = "debug"
	EnvEmpty = ""

	DebugConfigFile = "config/config.local.yaml"

	ConfigFileTemplate = "config/config.tmpl.yaml"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nest",
	Short: "`nest` command line tool",
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
	f.StringVarP(&cfgFile, "config", "c", "", "config file")
	f.String("env", "", "production, preview, debug")
	f.String("db-dialect", "sqlite", "database dialect [mysql, postgres, sqlite]")
	f.String("db-args", "", "database args")
	f.String("kube-config", magicconch.Getenv("KUBECONFIG", os.Getenv("HOME")+"/.kube/config"), "kube config file path")

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()

	env := viper.GetString("env")

	if cfgFile == "" && (env == EnvDebug || env == EnvEmpty) {
		fs := afero.NewOsFs()
		exist, err := afero.Exists(fs, DebugConfigFile)
		magicconch.Must(err)
		if !exist {
			f1, err := fs.Open(ConfigFileTemplate)
			magicconch.Must(err)
			defer func() {
				magicconch.Must(f1.Close())
			}()
			f2, err := fs.Create(DebugConfigFile)
			magicconch.Must(err)
			defer func() {
				magicconch.Must(f2.Close())
			}()
			_, err = io.Copy(f2, f1)
			magicconch.Must(err)
		}
		log.Infof("config file not specified, using default for debugging: %s", DebugConfigFile)
		cfgFile = DebugConfigFile
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
