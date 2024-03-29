package cache

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

const prefix = "heetch.driver-location"

type RedisDriver struct {
	client redis.Pool
}

func NewRedisDriver(address string) RedisDriver {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", address)
		},
		MaxIdle:   80,
		MaxActive: 12000,
	}

	return RedisDriver{client: pool}
}

type LocationDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at"` // date in RFC3330 format
}

type DriverLocationDTO struct {
	DriverID  string        `json:"driver_id"`
	Locations []LocationDTO `json:"locations"`
}

func (rd RedisDriver) Save(driver domain.Driver) error {
	conn := rd.client.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", prefix, driver.DriverID().String())
	driverLocationDTO := rd.entityToDTO(driver)
	bytes, err := json.Marshal(driverLocationDTO)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, bytes)

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
	var driverLocationDTO DriverLocationDTO
	err = json.Unmarshal(bytes, &driverLocationDTO)
	if err != nil {
		return domain.Driver{}, domain.ErrDriverNotFound
	}

	return rd.dtoToEntity(driverLocationDTO), nil
}

func (rd RedisDriver) entityToDTO(driver domain.Driver) DriverLocationDTO {
	locationsDTO := make([]LocationDTO, 0)
	for _, location := range driver.Locations() {
		locationDTO := LocationDTO{
			Latitude:  location.Latitude(),
			Longitude: location.Longitude(),
			UpdatedAt: location.UpdatedAt().String(),
		}
		locationsDTO = append(locationsDTO, locationDTO)
	}

	return DriverLocationDTO{
		DriverID:  driver.DriverID().String(),
		Locations: locationsDTO,
	}
}

func (rd RedisDriver) dtoToEntity(dto DriverLocationDTO) domain.Driver {

}
