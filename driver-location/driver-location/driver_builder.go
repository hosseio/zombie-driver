package driver_location

import (
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/internal"
)

type DriverBuilder struct{}

func NewDriverBuilder() DriverBuilder {
	return DriverBuilder{}
}

type LocationBuilderDTO struct {
	Latitude  float64
	Longitude float64
	UpdatedAt time.Time
}

type DriverBuilderDTO struct {
	DriverID  string
	Locations []LocationBuilderDTO
}

func (c DriverBuilder) Build(dto DriverBuilderDTO) (domain.Driver, error) {
	id, err := domain.NewDriverID(dto.DriverID)
	if err != nil {
		return domain.Driver{}, err
	}
	var locations domain.LocationList
	for _, locationDTO := range dto.Locations {
		location, err := c.BuildLocation(locationDTO)
		if err != nil {
			return domain.Driver{}, err
		}
		locations.Add(location)
	}

	return domain.NewDriver(id, locations)
}

func (c DriverBuilder) BuildLocation(dto LocationBuilderDTO) (domain.Location, error) {
	at, err := domain.NewUpdatedAt(dto.UpdatedAt)
	if err != nil {
		return domain.Location{}, err
	}

	return domain.NewLocation(dto.Latitude, dto.Longitude, at)
}
