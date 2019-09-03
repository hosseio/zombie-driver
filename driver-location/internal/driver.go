package domain

type Driver struct {
	driverID DriverID
	locations LocationList
}

func NewDriver(driverID DriverID) (Driver, error) {
	return Driver{driverID: driverID}, nil
}

func (d Driver) AddLocation(location Location)  {
	d.locations.add(location)
}
