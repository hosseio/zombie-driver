package driver_zombie

//go:generate moq -out distance_calculator_mock.go . DistanceCalculator
type DistanceCalculator interface {
	Calculate(driverID string, lastMinutes int) (int, error)
}
