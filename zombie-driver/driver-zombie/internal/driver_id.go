package domain

import (
	"errors"

	"github.com/satori/go.uuid"
)

type DriverID struct {
	value string
}

var ErrDriverIDEmpty = errors.New("driver id cannot be empty")

func NewDriverID(value string) (DriverID, error) {
	if value == "" {
		return DriverID{}, ErrDriverIDEmpty
	}
	_, err := uuid.FromString(value)
	if err != nil {
		return DriverID{}, err
	}

	return DriverID{value}, nil
}

func (driverID DriverID) String() string {
	return driverID.value
}

func (driverID DriverID) Equal(other DriverID) bool {
	return driverID.value == driverID.value
}
