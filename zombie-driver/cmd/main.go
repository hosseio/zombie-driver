package main

import (
	"context"

	"github.com/chiguirez/snout"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/bootstrap"
)

func main() {
	kernel := snout.Kernel{
		RunE: bootstrap.Run,
	}

	kernelBootstrap := kernel.Bootstrap(
		"zombie-driver",
		&bootstrap.Config{},
	)

	if err := kernelBootstrap.Initialize(); err != nil {
		if err != context.Canceled {
			panic(err)
		}
	}
}
