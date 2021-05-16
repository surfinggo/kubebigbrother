package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
	"k8s.io/client-go/util/homedir"
	"path"
)

type KubeconfigOptions struct {
	Kubeconfig string
}

func GetKubeconfigOptions() *KubeconfigOptions {
	return &KubeconfigOptions{
		Kubeconfig: viper.GetString("kubeconfig"),
	}
}

func AddKubeconfigFlags(fs *pflag.FlagSet) {
	defaultV := magicconch.Getenv("KUBECONFIG", path.Join(homedir.HomeDir(), ".kube", "config"))
	fs.String("kubeconfig", defaultV, "path to kubeconfig file")
}
