package driver_zombie

//go:generate moq -out is_zombie_resolver_mock.go . IsZombieResolver
type IsZombieResolver interface {
	Resolve(driverID string) bool
}

type DriverIsZombieResolver struct {
	distanceCalculator DistanceCalculator
	configGetter       ZombieConfigurer
}

func NewDriverIsZombieResolver(distanceCalculator DistanceCalculator, configGetter ZombieConfigurer) DriverIsZombieResolver {
	return DriverIsZombieResolver{distanceCalculator, configGetter}
}

func (r DriverIsZombieResolver) Resolve(driverID string) bool {
	config, err := r.configGetter.GetZombieConfig()
	if err != nil {
		return false
	}

	distance, err := r.distanceCalculator.Calculate(driverID, config.LastMinutes)
	if err != nil {
		return false
	}

	return distance < config.MaxMeters
}
