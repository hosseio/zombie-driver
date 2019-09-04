//+build integration

package cache

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location"

	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestRedisDriver(t *testing.T) {
	assertThat := require.New(t)
	redisAddr := os.Getenv("REDIS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	driverID := uuid.NewV4().String()

	cleanRedis := func() {
		conn, err := redis.Dial("tcp", redisAddr)
		if err != nil {
			println(err)
		}
		_, err = conn.Do("DEL", fmt.Sprintf("%s:%s", prefix, driverID))
		if err != nil {
			println(err)
		}
	}
	defer cleanRedis()
	t.Run("Given a redis driver", func(t *testing.T) {
		driverBuilder := driver_location.NewDriverBuilder()
		sut := NewRedisDriver(redisAddr, driverBuilder)
		t.Run("When saving a driver with locations", func(t *testing.T) {
			now := time.Now().UTC()
			dto := driver_location.DriverBuilderDTO{
				DriverID: driverID,
				Locations: []driver_location.LocationBuilderDTO{
					driver_location.LocationBuilderDTO{0.0, 0.0, now},
					driver_location.LocationBuilderDTO{1.0, 1.0, now},
				},
			}
			driver, _ := driverBuilder.Build(dto)
			err := sut.Save(driver)
			assertThat.NoError(err)
			t.Run("Then it can be retrieved", func(t *testing.T) {
				retrievedDriver, err := sut.ById(driverID)
				assertThat.NoError(err)
				assertThat.Equal(driver, retrievedDriver)
			})
		})
	})
}
