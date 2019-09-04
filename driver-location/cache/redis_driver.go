package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location"

	"github.com/gomodule/redigo/redis"
	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

const prefix = "heetch.driver-location"

type RedisDriver struct {
	client        redis.Pool
	driverBuilder driver_location.DriverBuilder
}

func NewRedisDriver(address string, driverBuilder driver_location.DriverBuilder) RedisDriver {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", address)
		},
		MaxIdle:   80,
		MaxActive: 12000,
	}

	return RedisDriver{client: pool, driverBuilder: driverBuilder}
}

type LocationDTO struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
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

	return rd.dtoToEntity(driverLocationDTO)
}

func (rd RedisDriver) entityToDTO(driver domain.Driver) DriverLocationDTO {
	locationsDTO := make([]LocationDTO, 0)
	for _, location := range driver.Locations() {
		locationDTO := LocationDTO{
			Latitude:  location.Latitude(),
			Longitude: location.Longitude(),
			UpdatedAt: location.UpdatedAt().Date(),
		}
		locationsDTO = append(locationsDTO, locationDTO)
	}

	return DriverLocationDTO{
		DriverID:  driver.DriverID().String(),
		Locations: locationsDTO,
	}
}

func (rd RedisDriver) dtoToEntity(dto DriverLocationDTO) (domain.Driver, error) {
	locations := []driver_location.LocationBuilderDTO{}
	for _, location := range dto.Locations {
		locations = append(locations, driver_location.LocationBuilderDTO{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
			UpdatedAt: location.UpdatedAt,
		})
	}

	builderDTO := driver_location.DriverBuilderDTO{
		DriverID:  dto.DriverID,
		Locations: locations,
	}

	return rd.driverBuilder.Build(builderDTO)
}
