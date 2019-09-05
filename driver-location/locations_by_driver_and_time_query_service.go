package driver_location

import (
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

type LocationsByDriverAndTimeQueryService struct {
	view        domain.LocationView
	transformer Transformer
}

func NewLocationsByDriverAndTimeQueryService(view domain.LocationView, transformer Transformer) LocationsByDriverAndTimeQueryService {
	return LocationsByDriverAndTimeQueryService{view: view, transformer: transformer}
}

func (qs LocationsByDriverAndTimeQueryService) Get(driverID string, from time.Time) ([]LocationDTO, error) {
	locations, err := qs.view.ByDriverAndFromDate(driverID, from)
	if err != nil {
		return []LocationDTO{}, err
	}

	return qs.transformer.LocationListToDTO(locations), nil
}
