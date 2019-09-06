package bootstrap

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg Config) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		server, err := InitializeServer(cfg)
		if err != nil {
			return err
		}

		return server.ListenAndServe()
	})

	g.Go(func() error {
		consumer, err := InitializeCreateDriverLocationNsqConsumer(cfg)
		if err != nil {
			return err
		}

		return consumer.Run(ctx)
	})

	return g.Wait()
}
