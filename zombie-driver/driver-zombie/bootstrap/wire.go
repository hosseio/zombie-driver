// +build wireinject

package bootstrap

import (
	netHttp "net/http"

	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/distance"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/zombie-configuration"

	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"

	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/http"

	"github.com/google/wire"
)

var HttpSet = wire.NewSet(
	http.NewServer,
	http.NewRouter,
	http.NewZombieController,
	http.NewConfigController,
)

var AppSet = wire.NewSet(
	driver_zombie.NewDriverIsZombieResolver,
	wire.Bind(new(driver_zombie.IsZombieResolver), driver_zombie.DriverIsZombieResolver{}),
)

var ZombieConfigurationSet = wire.NewSet(
	zombie_configuration.NewHardcodedZombieConfigGetter,
	zombie_configuration.NewRedisZombieConfigurer,
)

func InitializeRedisZombieConfigurer(cfg Config) (zombie_configuration.RedisZombieConfigurer, error) {
	wire.Build(
		ZombieConfigurationSet,
		getRedisAddr,
	)

	return zombie_configuration.RedisZombieConfigurer{}, nil
}

func getRedisAddr(cfg Config) zombie_configuration.RedisAddr {
	return zombie_configuration.RedisAddr(cfg.Redis)
}

func InitializeHardcodedZombieConfigGetter(cfg Config) (zombie_configuration.HardcodedZombieConfigGetter, error) {
	wire.Build(
		ZombieConfigurationSet,
	)

	return zombie_configuration.HardcodedZombieConfigGetter{}, nil
}

var DistanceSet = wire.NewSet(
	distance.NewLocationsDistanceCalculator,
	distance.NewDriverLocationClient,
	wire.Bind(new(distance.LocationsGetter), distance.DriverLocationClient{}),
)

func InitializeServer(cfg Config) (*netHttp.Server, error) {
	wire.Build(
		HttpSet,
		AppSet,
		serverAddr,
		InitializeLocationsDistanceCalculator,
		wire.Bind(new(driver_zombie.DistanceCalculator), distance.LocationsDistanceCalculator{}),
		InitializeRedisZombieConfigurer,
		wire.Bind(new(driver_zombie.ZombieConfigurer), zombie_configuration.RedisZombieConfigurer{}),
	)

	return &netHttp.Server{}, nil
}

func serverAddr(cfg Config) http.ServerAddr {
	return http.ServerAddr(cfg.Server.Addr)
}

func InitializeLocationsDistanceCalculator(cfg Config) (distance.LocationsDistanceCalculator, error) {
	wire.Build(
		DistanceSet,
		getBaseDriverLocationsURL,
	)

	return distance.LocationsDistanceCalculator{}, nil
}

func getBaseDriverLocationsURL(cfg Config) distance.BaseDriverLocationsURL {
	return distance.BaseDriverLocationsURL(cfg.DriverLocation.BaseURL)
}
