// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/distance"
	http2 "github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/http"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie/zombie-configuration"
	"net/http"
)

// Injectors from wire.go:

func InitializeHardcodedZombieConfigGetter(cfg Config) (zombie_configuration.HardcodedZombieConfigGetter, error) {
	hardcodedZombieConfigGetter := zombie_configuration.NewHardcodedZombieConfigGetter()
	return hardcodedZombieConfigGetter, nil
}

func InitializeServer(cfg Config) (*http.Server, error) {
	httpServerAddr := serverAddr(cfg)
	locationsDistanceCalculator, err := InitializeLocationsDistanceCalculator(cfg)
	if err != nil {
		return nil, err
	}
	hardcodedZombieConfigGetter, err := InitializeHardcodedZombieConfigGetter(cfg)
	if err != nil {
		return nil, err
	}
	driverIsZombieResolver := driver_zombie.NewDriverIsZombieResolver(locationsDistanceCalculator, hardcodedZombieConfigGetter)
	zombieController := http2.NewZombieController(driverIsZombieResolver)
	router := http2.NewRouter(zombieController)
	server := http2.NewServer(httpServerAddr, router)
	return server, nil
}

func InitializeLocationsDistanceCalculator(cfg Config) (distance.LocationsDistanceCalculator, error) {
	baseDriverLocationsURL := getBaseDriverLocationsURL(cfg)
	driverLocationClient := distance.NewDriverLocationClient(baseDriverLocationsURL)
	locationsDistanceCalculator := distance.NewLocationsDistanceCalculator(driverLocationClient)
	return locationsDistanceCalculator, nil
}

// wire.go:

var HttpSet = wire.NewSet(http2.NewServer, http2.NewRouter, http2.NewZombieController)

var AppSet = wire.NewSet(driver_zombie.NewDriverIsZombieResolver, wire.Bind(new(driver_zombie.IsZombieResolver), driver_zombie.DriverIsZombieResolver{}))

var ZombieConfigurationSet = wire.NewSet(zombie_configuration.NewHardcodedZombieConfigGetter)

var DistanceSet = wire.NewSet(distance.NewLocationsDistanceCalculator, distance.NewDriverLocationClient, wire.Bind(new(distance.LocationsGetter), distance.DriverLocationClient{}))

func serverAddr(cfg Config) http2.ServerAddr {
	return http2.ServerAddr(cfg.Server.Addr)
}

func getBaseDriverLocationsURL(cfg Config) distance.BaseDriverLocationsURL {
	return distance.BaseDriverLocationsURL(cfg.DriverLocation.BaseURL)
}
