package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewLocation(t *testing.T) {
	assertThat := require.New(t)
	now := time.Now()
	at, _ := NewUpdatedAt(now)
	t.Run("Given a latitude less than -90 value", func(t *testing.T) {
		lat := -100.0
		lon := 0.0
		t.Run("When the location is created", func(t *testing.T) {
			location, err := NewLocation(lat, lon, at)
			t.Run("Then an error occurred", func(t *testing.T) {
				assertThat.Error(err)
				assertThat.IsType(ErrInvalidLocation{}, err)
				assertThat.Empty(location)
			})
		})
	})
	t.Run("Given a longitude greater than 180 value", func(t *testing.T) {
		lat := 0.0
		lon := 190.0
		t.Run("When the location is created", func(t *testing.T) {
			location, err := NewLocation(lat, lon, at)
			t.Run("Then an error occurred", func(t *testing.T) {
				assertThat.Error(err)
				assertThat.IsType(ErrInvalidLocation{}, err)
				assertThat.Empty(location)
			})
		})
	})
	t.Run("Given latitude and longitude values in the permitted range", func(t *testing.T) {
		lat := 0.0
		lon := 0.0
		t.Run("When the location is created", func(t *testing.T) {
			location, err := NewLocation(lat, lon, at)
			t.Run("Then there is no error", func(t *testing.T) {
				assertThat.NoError(err)
			})
			t.Run("And the given values are set", func(t *testing.T) {
				assertThat.Equal(Location{lat, lon, at}, location)
			})
		})
	})
}
