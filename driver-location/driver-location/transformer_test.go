package driver_location

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestTransformer(t *testing.T) {
	assertThat := require.New(t)
	driverID := uuid.NewV4().String()

	t.Run("Given a driver with locations", func(t *testing.T) {
		now := time.Now()
		dto := DriverBuilderDTO{
			DriverID: driverID,
			Locations: []LocationBuilderDTO{
				{0.0, 0.0, now},
				{1.0, 1.0, now},
			},
		}
		builder := NewDriverBuilder()
		driver, _ := builder.Build(dto)
		t.Run("When is transformed to primitive data", func(t *testing.T) {
			sut := NewTransformer(builder)
			transformed := sut.DriverEntityToDTO(driver)
			t.Run("Then it is transformed properly", func(t *testing.T) {
				assertThat.Equal(driver.DriverID().String(), transformed.DriverID)
				assertThat.Equal(len(driver.Locations()), len(transformed.Locations))
				assertThat.Equal(driver.Locations()[0].Latitude(), transformed.Locations[0].Latitude)
			})
		})
	})
	t.Run("Given a DTO to be transformed", func(t *testing.T) {
		now := time.Now()
		plainDriver := DriverLocationDTO{
			DriverID: driverID,
			Locations: []LocationDTO{
				{0.0, 0.0, now},
				{1.0, 1.0, now},
			},
		}
		builder := NewDriverBuilder()
		t.Run("When is transformed to entity Driver", func(t *testing.T) {
			sut := NewTransformer(builder)
			dto := DriverBuilderDTO{
				DriverID: driverID,
				Locations: []LocationBuilderDTO{
					{0.0, 0.0, now},
					{1.0, 1.0, now},
				},
			}
			expected, _ := builder.Build(dto)
			transformed, err := sut.DTOToDriverEntity(plainDriver)
			t.Run("Then it is transformed properly", func(t *testing.T) {
				assertThat.NoError(err)
				assertThat.Equal(expected.DriverID().String(), transformed.DriverID().String())
				assertThat.Equal(len(expected.Locations()), len(transformed.Locations()))
				assertThat.Equal(expected.Locations()[0].Latitude(), transformed.Locations()[0].Latitude())
			})
		})
	})
}
