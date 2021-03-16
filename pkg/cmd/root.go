package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	"math/rand"
	"time"
)

var cfgFile string

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

	err := viper.BindPFlags(f)
	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Warn(errors.Wrapf(err, "read in config error, file: %s", viper.ConfigFileUsed()))
		} else {
			log.Infof("using config file: %s", viper.ConfigFileUsed())
		}
	}
}
