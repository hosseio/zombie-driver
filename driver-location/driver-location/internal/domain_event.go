package domain

import "time"

type DriverCreated struct {
	DriverID  string        `json:"driver_id"`
	Locations []LocationDTO `json:"locations"`
}

func (e DriverCreated) AggregateID() string {
	return e.DriverID
}

type LocationAddedToDriver struct {
	DriverID  string    `json:"driver_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e LocationAddedToDriver) AggregateID() string {
	return e.DriverID
}

type LocationDTO struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DriverLocationDTO struct {
	DriverID  string        `json:"driver_id"`
	Locations []LocationDTO `json:"locations"`
}
