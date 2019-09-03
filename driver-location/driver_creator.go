package driver_location

import "github.com/heetch/jose-odg-technical-test/driver-location/internal"

type DriverCreator struct{}

func NewDriverCreator() DriverCreator {
	return DriverCreator{}
}

func (c DriverCreator) Create(driverID string) (domain.Driver, error) {
	id, err := domain.NewDriverID(driverID)
	if err != nil {
		return domain.Driver{}, err
	}

	return domain.NewDriver(id)
}
