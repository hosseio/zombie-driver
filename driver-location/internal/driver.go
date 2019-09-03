package domain

type Driver struct {
	driverID  DriverID
	locations LocationList
}

func NewDriver(driverID DriverID, locations LocationList) (Driver, error) {
	return Driver{driverID, locations}, nil
}

func (d Driver) AddLocation(location Location) {
	d.locations.Add(location)
}

func (d Driver) Locations() LocationList {
	return d.locations
}

func (d Driver) DriverID() DriverID {
	return d.driverID
}
