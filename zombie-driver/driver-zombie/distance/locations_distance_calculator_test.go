package distance

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/require"
)

func TestLocationsDistanceCalculator_Calculate(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given a locations distance calculator", func(t *testing.T) {
		locationsGetterMock := &LocationsGetterMock{}
		sut := NewLocationsDistanceCalculator(locationsGetterMock)
		t.Run("When two locations are given", func(t *testing.T) {
			locationsGetterMock.GetLocationsFunc = func(driverID string, lastMinutes int) (lists LocationList, e error) {
				locationList := LocationList{}
				locationList.Add(Location{48.864193, 2.350498, time.Now()})
				locationList.Add(Location{48.863921, 2.349211, time.Now()})

				return locationList, nil
			}
			t.Run("Then the distance is properly calculated", func(t *testing.T) {
				driverID := uuid.NewV4().String()
				distance, err := sut.Calculate(driverID, 5)
				assertThat.NoError(err)
				assertThat.Equal(98, distance)
			})
		})
		t.Run("When three locations are given", func(t *testing.T) {
			locationsGetterMock.GetLocationsFunc = func(driverID string, lastMinutes int) (lists LocationList, e error) {
				locationList := LocationList{}
				locationList.Add(Location{48.864193, 2.350498, time.Now()})
				locationList.Add(Location{48.863921, 2.349211, time.Now()})
				locationList.Add(Location{48.864193, 2.350498, time.Now()}) // back to origin

				return locationList, nil
			}
			t.Run("Then the distance is properly calculated", func(t *testing.T) {
				driverID := uuid.NewV4().String()
				distance, err := sut.Calculate(driverID, 5)
				assertThat.NoError(err)
				assertThat.Equal(197, distance)
			})
		})
		t.Run("When one location is given", func(t *testing.T) {
			locationsGetterMock.GetLocationsFunc = func(driverID string, lastMinutes int) (lists LocationList, e error) {
				locationList := LocationList{}
				locationList.Add(Location{48.864193, 2.350498, time.Now()})

				return locationList, nil
			}
			t.Run("Then the distance is zero", func(t *testing.T) {
				driverID := uuid.NewV4().String()
				distance, err := sut.Calculate(driverID, 5)
				assertThat.NoError(err)
				assertThat.Equal(0, distance)
			})
		})
	})
}
