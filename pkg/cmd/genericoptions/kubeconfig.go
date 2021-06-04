package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// KubeconfigOptions is kubeconfig options
type KubeconfigOptions struct {
	Kubeconfig string
}

// GetKubeconfigOptions get kubeconfig options from viper flags
func GetKubeconfigOptions() *KubeconfigOptions {
	return &KubeconfigOptions{
		Kubeconfig: viper.GetString("kubeconfig"),
	}
}

// AddKubeconfigFlags adds kubeconfig flags to flag set
func AddKubeconfigFlags(fs *pflag.FlagSet) {
	fs.String("kubeconfig", os.Getenv("KUBECONFIG"), "path to kubeconfig file")
}
