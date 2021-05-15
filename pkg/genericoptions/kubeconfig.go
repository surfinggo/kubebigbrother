package genericoptions

import (
	"github.com/spf13/pflag"
	"github.com/spongeprojects/kubebigbrother/pkg/helpers/homedir"
	"github.com/spongeprojects/magicconch"
	"path"
)

type KubeconfigOptions struct {
	Kubeconfig string
}

func GetKubeconfigOptions() *KubeconfigOptions {
	return &KubeconfigOptions{
		Kubeconfig: magicconch.Getenv("KUBECONFIG", path.Join(homedir.HomeDir(), ".kube", "config")),
	}
}

func AddKubeconfigFlags(fs *pflag.FlagSet) {
	defaultV := magicconch.Getenv("KUBECONFIG", path.Join(homedir.HomeDir(), ".kube", "config"))
	fs.String("kubeconfig", defaultV, "path to kubeconfig file")
}
