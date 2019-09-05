package cache

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location"

	"github.com/gomodule/redigo/redis"
	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

const (
	prefix     = "heetch.driver-location"
	ttlInHours = 48
)

type RedisDriver struct {
	client      redis.Pool
	transformer driver_location.Transformer
}

func NewRedisDriver(address string, transformer driver_location.Transformer) RedisDriver {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", address)
		},
		MaxIdle:   80,
		MaxActive: 12000,
	}

	return RedisDriver{client: pool, transformer: transformer}
}

func (rd RedisDriver) Save(driver domain.Driver) error {
	conn := rd.client.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", prefix, driver.DriverID().String())

	driverLocationDTO := rd.transformer.DriverEntityToDTO(driver)
	bytes, err := json.Marshal(driverLocationDTO)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", key, ttlInHours*3600, bytes)

	return err
}

func (rd RedisDriver) ById(id string) (domain.Driver, error) {
	conn := rd.client.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", prefix, id)

	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return domain.Driver{}, domain.ErrDriverNotFound
	}

	var driverLocationDTO driver_location.DriverLocationDTO
	err = json.Unmarshal(bytes, &driverLocationDTO)
	if err != nil {
		return domain.Driver{}, domain.ErrDriverNotFound
	}

	return rd.transformer.DTOToDriverEntity(driverLocationDTO)
}

func (rd RedisDriver) ByDriverAndFromDate(id string, from time.Time) (domain.LocationList, error) {
	driver, err := rd.ById(id)
	if err != nil {
		return domain.LocationList{}, err
	}

	locations := domain.LocationList{}
	for _, location := range driver.Locations() {
		if location.UpdatedAt().IsAfter(from) {
			locations.Add(location)
		}
	}

	sort.Sort(driver_location.ByUpdatedAt(locations))

	return locations, nil
}
