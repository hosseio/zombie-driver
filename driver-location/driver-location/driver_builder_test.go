package driver_location

import (
	"testing"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/internal"
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/require"
)

func TestNewDriverBuilder(t *testing.T) {
	assertThat := require.New(t)
	t.Run("Given a builder and a dto", func(t *testing.T) {
		now := time.Now()
		dto := DriverBuilderDTO{
			DriverID: uuid.NewV4().String(),
			Locations: []LocationBuilderDTO{
				LocationBuilderDTO{0.0, 0.0, now},
				LocationBuilderDTO{1.0, 1.0, now},
			},
		}
		sut := NewDriverBuilder()
		t.Run("When the builder is executed", func(t *testing.T) {
			driver, err := sut.Build(dto)
			t.Run("Then a driver is returned", func(t *testing.T) {
				assertThat.NoError(err)
				assertThat.IsType(domain.Driver{}, driver)
				assertThat.Equal(dto.DriverID, driver.DriverID().String())
				assertThat.Equal(2, len(driver.Locations()))
			})
		})
	})
}
