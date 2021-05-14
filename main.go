package main

import (
	"github.com/spongeprojects/kubebigbrother/pkg/cmd"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cmd.Execute()
}
