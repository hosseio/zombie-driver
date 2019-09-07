package zombie_configuration

import (
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"
)

type HardcodedZombieConfigGetter struct{}

func NewHardcodedZombieConfigGetter() HardcodedZombieConfigGetter {
	return HardcodedZombieConfigGetter{}
}

const MaxMeters = 500
const LastMinutes = 5

func (HardcodedZombieConfigGetter) GetZombieConfig() driver_zombie.ZombieConfigProjection {
	return driver_zombie.ZombieConfigProjection{
		MaxMeters:   MaxMeters,
		LastMinutes: LastMinutes,
	}
}
