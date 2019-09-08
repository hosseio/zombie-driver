package main

import (
	"context"

	"github.com/chiguirez/snout"
	"github.com/heetch/jose-odg-technical-test/gateway/bootstrap"
)

const urlFilename = "config.yaml"

func main() {
	kernel := snout.Kernel{
		RunE: bootstrap.Run,
	}

	config := bootstrap.Config{}
	kernelBootstrap := kernel.Bootstrap(
		"gateway",
		&config,
	)

	if err := config.ReadURLConfiguration(urlFilename); err != nil {
		panic(err)
	}

	if err := kernelBootstrap.Initialize(); err != nil {
		if err != context.Canceled {
			panic(err)
		}
	}
}
