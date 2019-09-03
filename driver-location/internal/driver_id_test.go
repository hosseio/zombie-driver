package domain

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewDriverID(t *testing.T) {
	assertThat := require.New(t)
	t.Run("Given a valid uuid", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		t.Run("When the driverID is created", func(t *testing.T) {
			ID, err := NewDriverID(uuid)
			t.Run("Then it is created with that value and without errors", func(t *testing.T) {
				assertThat.NoError(err)
				assertThat.Equal(DriverID{uuid}, ID)
			})
		})
	})
	t.Run("Given a NON valid uuid", func(t *testing.T) {
		uuid := "something"
		t.Run("When the driverID is created", func(t *testing.T) {
			_, err := NewDriverID(uuid)
			t.Run("Then an error is returned", func(t *testing.T) {
				assertThat.Error(err)
			})
		})
	})
	t.Run("Given an empty value", func(t *testing.T) {
		uuid := ""
		t.Run("When the driverID is created", func(t *testing.T) {
			_, err := NewDriverID(uuid)
			t.Run("Then an error is returned", func(t *testing.T) {
				assertThat.Error(err)
			})
		})
	})
}

func TestDriverID_Equal(t *testing.T) {
	assertThat := require.New(t)
	t.Run("Given two DriverIDs with the same value", func(t *testing.T) {
		uuid := uuid.NewV4().String()
		v1, _ := NewDriverID(uuid)
		v2, _ := NewDriverID(uuid)
		t.Run("When checking if they are equal", func(t *testing.T) {
			equal := v1.Equal(v2)
			t.Run("Then a true is returned", func(t *testing.T) {
				assertThat.True(equal)
			})
		})
	})
	t.Run("Given two DriverIDs with different value", func(t *testing.T) {
		v1, _ := NewDriverID(uuid.NewV4().String())
		v2, _ := NewDriverID(uuid.NewV4().String())
		t.Run("When checking if they are equal", func(t *testing.T) {
			equal := v1.Equal(v2)
			t.Run("Then a false is returned", func(t *testing.T) {
				assertThat.True(equal)
			})
		})
	})
}
