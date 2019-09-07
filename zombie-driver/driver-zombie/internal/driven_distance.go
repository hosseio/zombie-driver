package domain

type DrivenDistance struct {
	driverID DriverID
	meters   int
	minutes  int
}

func NewDrivenDistance(driverID DriverID, meters int, minutes int) DrivenDistance {
	return DrivenDistance{driverID, meters, minutes}
}

func (d DrivenDistance) Minutes() int {
	return d.minutes
}

func (d DrivenDistance) Meters() int {
	return d.meters
}

func (d DrivenDistance) DriverID() DriverID {
	return d.driverID
}
