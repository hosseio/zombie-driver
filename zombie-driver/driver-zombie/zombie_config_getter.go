package driver_zombie

type ZombieConfigProjection struct {
	MaxMeters   int
	LastMinutes int
}

//go:generate moq -out zombie_config_getter_mock.go . ZombieConfigGetter
type ZombieConfigGetter interface {
	GetZombieConfig() ZombieConfigProjection
}
