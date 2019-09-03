package driver_location

import (
	"errors"
	"time"

	"github.com/chiguirez/cromberbus"
	"github.com/heetch/jose-odg-technical-test/driver-location/internal"
)

var ErrNonCreateLocationCommand = errors.New("Cannot handle a non create location command")

type CreateLocationCommand struct {
	Latitude  float64
	Longitude float64
	DriverID  string
}

type CreateLocationCommandHandler struct {
	driverView       domain.DriverView
	driverRepository domain.DriverRepository
	driverCreator    DriverCreator
}

func NewCreateLocationCommandHandler(
	driverView domain.DriverView,
	driverRepository domain.DriverRepository,
	driverCreator DriverCreator,
) CreateLocationCommandHandler {
	return CreateLocationCommandHandler{driverView: driverView, driverRepository: driverRepository, driverCreator: driverCreator}
}

func (h CreateLocationCommandHandler) Handle(command cromberbus.Command) error {
	createLocationCommand, ok := command.(CreateLocationCommand)
	if !ok {
		return ErrNonCreateLocationCommand
	}

	driver, err := h.getDriver(createLocationCommand.DriverID)
	if err != nil {
		return err
	}
	at, err := domain.NewUpdatedAt(time.Now())
	if err != nil {
		return err
	}
	location, err := domain.NewLocation(createLocationCommand.Latitude, createLocationCommand.Longitude, at)
	if err != nil {
		return err
	}

	driver.AddLocation(location)

	return h.driverRepository.Save(driver)
}

func (h CreateLocationCommandHandler) getDriver(driverID string) (domain.Driver, error) {
	driver, err := h.driverView.ById(driverID)
	if err != nil && err == domain.ErrDriverNotFound {
		return h.driverCreator.Create(driverID)
	}

	return driver, err
}
