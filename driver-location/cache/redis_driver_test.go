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
			yesterday := now.AddDate(0, 0, -1)
			yesterdayAtStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 00, 00, 00, 00, yesterday.Location())
			dto := driver_location.DriverBuilderDTO{
				DriverID: driverID,
				Locations: []driver_location.LocationBuilderDTO{
					{1.0, 1.0, yesterdayAtStart},
					{2.0, 2.0, now},
					{0.0, 0.0, yesterday},
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
			t.Run("And the locations can be retrieved filtering by date sorted", func(t *testing.T) {
				locations, err := sut.ByDriverAndFromDate(driverID, yesterday)
				assertThat.NoError(err)
				assertThat.Equal(2, len(locations))
				assertThat.Equal(0.0, locations[0].Latitude())
				assertThat.Equal(2.0, locations[1].Latitude())
			})
		})
	})
}
