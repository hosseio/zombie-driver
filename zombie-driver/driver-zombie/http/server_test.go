package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"
	"github.com/satori/go.uuid"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	assertThat := require.New(t)

	var (
		router             *mux.Router
		zombieController   ZombieController
		zombieResolverMock *driver_zombie.IsZombieResolverMock
	)

	setup := func() {
		zombieResolverMock = &driver_zombie.IsZombieResolverMock{}
		zombieController = NewLocationController(zombieResolverMock)
		router = NewRouter(zombieController)
	}

	t.Run("Given a server", func(t *testing.T) {
		setup()
		driverID := uuid.NewV4().String()
		zombieResolverMock.ResolveFunc = func(driverID string) bool {
			return true
		}
		t.Run("When consuming is zombie endpoint", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/drivers/"+driverID, nil)
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			t.Run("Then the locations are retrieved", func(t *testing.T) {
				body, err := ioutil.ReadAll(res.Body)
				assertThat.NoError(err)

				var zombieResponse ZombieResponse
				err = json.Unmarshal(body, &zombieResponse)
				assertThat.NoError(err)

				assertThat.Equal(http.StatusOK, res.Code)
				assertThat.Equal(driverID, zombieResponse.Id)
				assertThat.True(zombieResponse.Zombie)
			})
		})
	})
}
