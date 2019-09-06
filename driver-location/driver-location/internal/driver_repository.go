package domain

//go:generate moq -out driver_repository_mock.go . DriverRepository
type DriverRepository interface {
	Save(Driver) error
}
