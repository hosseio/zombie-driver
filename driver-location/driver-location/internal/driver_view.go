package domain

import "errors"

var ErrDriverNotFound = errors.New("Driver not found")

//go:generate moq -out driver_view_mock.go . DriverView
type DriverView interface {
	ById(string) (Driver, error)
}
