package driver_location

import (
	"sort"
	"testing"
	"time"

	"github.com/heetch/jose-odg-technical-test/driver-location/internal"

	"github.com/stretchr/testify/require"
)

func TestSortLocationsByUpdatedAt(t *testing.T) {
	assertThat := require.New(t)
	driverBuilder := NewDriverBuilder()

	t.Run("Given a location list unsorted", func(t *testing.T) {
		locations := domain.LocationList{}
		now := time.Now()
		l4, _ := driverBuilder.BuildLocation(LocationBuilderDTO{4.4, 4.4, now})
		l1, _ := driverBuilder.BuildLocation(LocationBuilderDTO{1.1, 1.1, now.AddDate(0, 0, -4)})
		l3, _ := driverBuilder.BuildLocation(LocationBuilderDTO{3.3, 3.3, now.AddDate(0, 0, -1)})
		l2, _ := driverBuilder.BuildLocation(LocationBuilderDTO{2.2, 2.2, now.AddDate(0, 0, -2)})

		locations.Add(l4)
		locations.Add(l1)
		locations.Add(l3)
		locations.Add(l2)
		t.Run("When it is sorted by updated at", func(t *testing.T) {
			sort.Sort(ByUpdatedAt(locations))
			t.Run("Then it is sorted", func(t *testing.T) {
				assertThat.Equal(l1, locations[0])
				assertThat.Equal(l2, locations[1])
				assertThat.Equal(l3, locations[2])
				assertThat.Equal(l4, locations[3])
			})
		})
	})
}
