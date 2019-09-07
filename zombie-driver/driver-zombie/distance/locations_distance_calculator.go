package distance

import (
	"math"
	"time"
)

type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocationList []Location

func (l *LocationList) Add(location Location) {
	*l = append([]Location(*l), location)
}

func (l LocationList) Len() int {
	return len(l)
}

func (l LocationList) Get(index int) Location {
	return l[index]
}

type LocationsDistanceCalculator struct {
	locationsGetter LocationsGetter
}

func NewLocationsDistanceCalculator(locationsGetter LocationsGetter) LocationsDistanceCalculator {
	return LocationsDistanceCalculator{locationsGetter: locationsGetter}
}

func (c LocationsDistanceCalculator) Calculate(driverID string, lastMinutes int) (int, error) {
	locations, err := c.locationsGetter.GetLocations(driverID, lastMinutes)
	if err != nil {
		return 0, err
	}

	return c.calculateDistanceFromLocations(locations), nil
}

func (c LocationsDistanceCalculator) calculateDistanceFromLocations(locations LocationList) int {
	if locations.Len() < 2 {
		return 0
	}

	distance := 0.0
	for i := 0; i < locations.Len()-1; i++ {
		startLocation := locations.Get(i)
		endLocation := locations.Get(i + 1)

		distance += c.calculateDistanceBetweenLocations(startLocation, endLocation)
	}

	return int(distance)
}

// Adapted from https://www.geodatasource.com/developers/go
func (c LocationsDistanceCalculator) calculateDistanceBetweenLocations(startLocation, endLocation Location) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * startLocation.Latitude / 180)
	radlat2 := float64(PI * endLocation.Latitude / 180)

	theta := float64(startLocation.Longitude - endLocation.Longitude)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	//km * 1000m
	dist = dist * 1.609344 * 1000

	return dist
}
