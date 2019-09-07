package bootstrap

import (
	"context"
)

func Run(ctx context.Context, cfg Config) error {
	server, err := InitializeServer(cfg)
	if err != nil {
		return err
	}

	return server.ListenAndServe()
}
