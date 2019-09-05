package driver_location

import (
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

type Transformer struct {
	driverBuilder DriverBuilder
}

func NewTransformer(driverBuilder DriverBuilder) Transformer {
	return Transformer{driverBuilder}
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

func (t Transformer) DriverEntityToDTO(driver domain.Driver) DriverLocationDTO {
	locationsDTO := t.LocationListToDTO(driver.Locations())

	return DriverLocationDTO{
		DriverID:  driver.DriverID().String(),
		Locations: locationsDTO,
	}
}

func (t Transformer) LocationListToDTO(locations domain.LocationList) []LocationDTO {
	locationsDTO := make([]LocationDTO, 0)
	for _, location := range locations {
		locationsDTO = append(locationsDTO, t.LocationEntityToDTO(location))
	}

	return locationsDTO
}

func (t Transformer) LocationEntityToDTO(location domain.Location) LocationDTO {
	return LocationDTO{
		Latitude:  location.Latitude(),
		Longitude: location.Longitude(),
		UpdatedAt: location.UpdatedAt().Date(),
	}
}

func (t Transformer) DTOToDriverEntity(dto DriverLocationDTO) (domain.Driver, error) {
	locations := []LocationBuilderDTO{}
	for _, location := range dto.Locations {
		locations = append(locations, LocationBuilderDTO{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
			UpdatedAt: location.UpdatedAt,
		})
	}

	builderDTO := DriverBuilderDTO{
		DriverID:  dto.DriverID,
		Locations: locations,
	}

	return t.driverBuilder.Build(builderDTO)
}
