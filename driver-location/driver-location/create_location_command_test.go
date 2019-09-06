package driver_location

import (
	"testing"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/internal"
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/require"
)

func TestCreateLocationCommandHandler_Handle(t *testing.T) {
	assertThat := require.New(t)
	t.Run("Given a create location command handler", func(t *testing.T) {
		driverViewMock := &domain.DriverViewMock{
			ByIdFunc: func(id string) (driver domain.Driver, e error) {
				return domain.Driver{}, nil
			},
		}
		driverRepositoryMock := &domain.DriverRepositoryMock{
			SaveFunc: func(d domain.Driver) error {
				return nil
			},
		}
		driverCreator := NewDriverBuilder()
		sut := NewCreateLocationCommandHandler(driverViewMock, driverRepositoryMock, driverCreator)
		t.Run("When it handles the create location command", func(t *testing.T) {
			command := CreateLocationCommand{
				DriverID:  uuid.NewV4().String(),
				Longitude: 0.0,
				Latitude:  0.0,
			}
			err := sut.Handle(command)
			t.Run("Then the driver of the new location is saved", func(t *testing.T) {
				assertThat.NoError(err)
				assertThat.True(len(driverRepositoryMock.SaveCalls()) > 0)
			})
		})
		t.Run("When it handles a non create location command", func(t *testing.T) {
			command := struct{}{}
			err := sut.Handle(command)
			t.Run("Then an error is returned", func(t *testing.T) {
				assertThat.Equal(ErrNonCreateLocationCommand, err)
			})
		})
	})
}
