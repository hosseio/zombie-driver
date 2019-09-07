package driver_zombie

//go:generate moq -out is_zombie_resolver_mock.go . IsZombieResolver
type IsZombieResolver interface {
	Resolve(driverID string) bool
}

type DriverIsZombieResolver struct {
	distanceCalculator DistanceCalculator
	configGetter       ZombieConfigGetter
}

func NewDriverIsZombieResolver(distanceCalculator DistanceCalculator, configGetter ZombieConfigGetter) DriverIsZombieResolver {
	return DriverIsZombieResolver{distanceCalculator, configGetter}
}

func (r DriverIsZombieResolver) Resolve(driverID string) bool {
	config := r.configGetter.GetZombieConfig()

	distance := r.distanceCalculator.Calculate(driverID, config.LastMinutes)

	return distance < config.MaxMeters
}
