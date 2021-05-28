package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
	"k8s.io/client-go/util/homedir"
	"path"
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
	defaultV := magicconch.Getenv("KUBECONFIG", path.Join(homedir.HomeDir(), ".kube", "config"))
	fs.String("kubeconfig", defaultV, "path to kubeconfig file")
}
