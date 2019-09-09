package domain

import "github.com/heetch/jose-odg-technical-test/driver-location/pkg"

type Driver struct {
	pkg.BaseAggregateRoot
	driverID  DriverID
	locations LocationList
}

func NewDriver(driverID DriverID, locations LocationList) (Driver, error) {
	driver := Driver{
		driverID:  driverID,
		locations: locations,
	}

	driver.Record(DriverCreated{
		driver.driverID.String(),
		driver.locations.toEventDTO(),
	})

	return driver, nil
}

func (d Driver) AddLocation(location Location) {
	d.locations.Add(location)

	d.Record(LocationAddedToDriver{
		DriverID:  d.driverID.String(),
		Latitude:  location.latitude,
		Longitude: location.longitude,
		UpdatedAt: location.updatedAt.Date(),
	})
}

func (d Driver) Locations() LocationList {
	return d.locations
}

func (d Driver) DriverID() DriverID {
	return d.driverID
}
