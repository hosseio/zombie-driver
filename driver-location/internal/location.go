package domain

import (
	"fmt"
)

type Location struct {
	latitude  float64
	longitude float64
	updatedAt UpdatedAt
}

type LocationList []Location

func (l *LocationList) Add(location Location) {
	*l = append([]Location(*l), location)
}

const (
	minLat = -90
	maxLat = 90
	minLon = -180
	maxLon = 180
)

type ErrInvalidLocation struct {
	message string
}

func (e ErrInvalidLocation) Error() string {
	return e.message
}

func NewLocation(lat float64, lon float64, at UpdatedAt) (Location, error) {
	if lat < minLat || lat > maxLat {
		return Location{}, ErrInvalidLocation{"invalid latitude"}
	}
	if lon < minLon || lon > maxLon {
		return Location{}, ErrInvalidLocation{"invalid longitude"}
	}

	return Location{lat, lon, at}, nil
}

func (l Location) String() string {
	return fmt.Sprintf("%v,%v", l.latitude, l.longitude)
}

func (l Location) Latitude() float64 {
	return l.latitude
}

func (l Location) Longitude() float64 {
	return l.longitude
}

func (l Location) UpdatedAt() UpdatedAt {
	return l.updatedAt
}
