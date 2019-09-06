// +build wireinject

package bootstrap

import (
	netHttp "net/http"

	"github.com/heetch/jose-odg-technical-test/driver-location/internal"

	"github.com/google/wire"
	"github.com/heetch/jose-odg-technical-test/driver-location"
	"github.com/heetch/jose-odg-technical-test/driver-location/cache"
	"github.com/heetch/jose-odg-technical-test/driver-location/http"
)

var ServerSet = wire.NewSet(
	http.NewServer,
	http.NewRouter,
	http.NewLocationController,

	InitializeRedisDriver,
	wire.Bind(new(domain.LocationView), cache.RedisDriver{}),
)

var RedisSet = wire.NewSet(
	driver_location.NewDriverBuilder,
)

var AppSet = wire.NewSet(
	driver_location.NewLocationsByDriverAndTimeQueryService,
	driver_location.NewTransformer,
)

func InitializeServer(cfg Config) (*netHttp.Server, error) {
	wire.Build(
		ServerSet,
		AppSet,
		RedisSet,
		serverAddr,
	)

	return &netHttp.Server{}, nil
}

func serverAddr(cfg Config) http.ServerAddr {
	return http.ServerAddr(cfg.Server.Addr)
}

func InitializeRedisDriver(cfg Config) (cache.RedisDriver, error) {
	wire.Build(
		RedisSet,
		AppSet,
		redisAddr,
		cache.NewRedisDriver,
	)

	return cache.RedisDriver{}, nil
}

func redisAddr(cfg Config) cache.RedisAddr {
	return cache.RedisAddr(cfg.Redis)
}
