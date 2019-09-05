package domain

import (
	"time"
)

//go:generate moq -out location_view_mock.go . LocationView
type LocationView interface {
	ByDriverAndFromDate(string, time.Time) (LocationList, error)
}
