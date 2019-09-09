package driver_zombie

type ZombieConfigProjection struct {
	MaxMeters   int `json:"max_meters"`
	LastMinutes int `json:"last_minutes"`
}

//go:generate moq -out zombie_configurer_mock.go . ZombieConfigurer
type ZombieConfigurer interface {
	GetZombieConfig() (ZombieConfigProjection, error)
	SetZombieConfig(ZombieConfigProjection) error
}
