package driver_location

import (
	"testing"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location/internal"
	"github.com/stretchr/testify/require"
)

func TestLocationsByDriverAndTimeQueryService(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given a locations by driver and time query service", func(t *testing.T) {
		builder := NewDriverBuilder()
		location, _ := builder.BuildLocation(LocationBuilderDTO{0.0, 0.0, time.Now()})
		view := &domain.LocationViewMock{
			ByDriverAndFromDateFunc: func(in1 string, in2 time.Time) (domain.LocationList, error) {
				return domain.LocationList{location}, nil
			},
		}
		// TODO mock dependencies
		sut := NewLocationsByDriverAndTimeQueryService(view, NewTransformer(builder))
		t.Run("When the data is asked", func(t *testing.T) {
			result, err := sut.Get("id", time.Now())
			t.Run("Then it returns the entities transformed", func(t *testing.T) {
				assertThat.NoError(err)
				assertThat.Equal(location.Latitude(), result[0].Latitude)
			})
		})
	})
}
