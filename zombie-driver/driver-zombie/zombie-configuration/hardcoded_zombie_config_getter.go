package zombie_configuration

import (
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"
)

type HardcodedZombieConfigGetter struct{}

const MaxMeters = 500
const LastMinutes = 5

func NewHardcodedZombieConfigGetter() HardcodedZombieConfigGetter {
	return HardcodedZombieConfigGetter{}
}

func (HardcodedZombieConfigGetter) GetZombieConfig() (driver_zombie.ZombieConfigProjection, error) {
	return driver_zombie.ZombieConfigProjection{
		MaxMeters:   MaxMeters,
		LastMinutes: LastMinutes,
	}, nil
}

func (HardcodedZombieConfigGetter) SetZombieConfig(driver_zombie.ZombieConfigProjection) error {
	return nil
}
