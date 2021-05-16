package main

import (
	"github.com/spongeprojects/kubebigbrother/pkg/cmd"
	"k8s.io/klog/v2"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	defer klog.Flush()

	command := cmd.NewKbbCommand()

	if err := command.Execute(); err != nil {
		klog.Exit(err)
	}
}
