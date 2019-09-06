package main

import (
	"context"

	bootstrap "github.com/heetch/jose-odg-technical-test/driver-location/driver-location/bootstrap"

	"github.com/chiguirez/snout"
)

func main() {
	kernel := snout.Kernel{
		RunE: bootstrap.Run,
	}

	kernelBootstrap := kernel.Bootstrap(
		"driver-location",
		&bootstrap.Config{},
	)

	if err := kernelBootstrap.Initialize(); err != nil {
		if err != context.Canceled {
			panic(err)
		}
	}
}
