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

	driverLocationDTO := rd.driverEntityToDTO(driver)
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

	var driverLocationDTO DriverLocationDTO
	err = json.Unmarshal(bytes, &driverLocationDTO)
	if err != nil {
		return domain.Driver{}, domain.ErrDriverNotFound
	}

	return rd.dtoToDriverEntity(driverLocationDTO)
}

func (rd RedisDriver) driverEntityToDTO(driver domain.Driver) DriverLocationDTO {
	locationsDTO := make([]LocationDTO, 0)
	for _, location := range driver.Locations() {
		locationsDTO = append(locationsDTO, rd.locationEntityToDTO(location))
	}

	return DriverLocationDTO{
		DriverID:  driver.DriverID().String(),
		Locations: locationsDTO,
	}
}

func (rd RedisDriver) locationEntityToDTO(location domain.Location) LocationDTO {
	return LocationDTO{
		Latitude:  location.Latitude(),
		Longitude: location.Longitude(),
		UpdatedAt: location.UpdatedAt().Date(),
	}
}

func (rd RedisDriver) dtoToDriverEntity(dto DriverLocationDTO) (domain.Driver, error) {
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
