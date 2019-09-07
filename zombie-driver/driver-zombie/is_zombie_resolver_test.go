package driver_zombie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriverIsZombieResolver_Resolve(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given a driver is zombie resolver service", func(t *testing.T) {
		distanceCalculatorMock := &DistanceCalculatorMock{}
		zombieConfigGetterMock := &ZombieConfigGetterMock{
			GetZombieConfigFunc: func() ZombieConfigProjection {
				return ZombieConfigProjection{
					MaxMeters:   500,
					LastMinutes: 5,
				}
			},
		}
		sut := NewDriverIsZombieResolver(distanceCalculatorMock, zombieConfigGetterMock)
		t.Run("When the calculated distance for the given config is less than the config max meters", func(t *testing.T) {
			distanceCalculatorMock.CalculateFunc = func(driverID string, lastMinutes int) int {
				return 200
			}
			t.Run("Then the driver is zombie", func(t *testing.T) {
				assertThat.True(sut.Resolve("driver-id"))
			})
		})
		t.Run("When the calculated distance for the given config is equal than the config max meters", func(t *testing.T) {
			distanceCalculatorMock.CalculateFunc = func(driverID string, lastMinutes int) int {
				return 500
			}
			t.Run("Then the driver is NOT zombie", func(t *testing.T) {
				assertThat.False(sut.Resolve("driver-id"))
			})
		})
		t.Run("When the calculated distance for the given config is greater than the config max meters", func(t *testing.T) {
			distanceCalculatorMock.CalculateFunc = func(driverID string, lastMinutes int) int {
				return 501
			}
			t.Run("Then the driver is NOT zombie", func(t *testing.T) {
				assertThat.False(sut.Resolve("driver-id"))
			})
		})
	})
}
