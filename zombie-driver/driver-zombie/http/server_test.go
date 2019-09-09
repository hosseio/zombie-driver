package http

import (
	"bytes"
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
		sut                *mux.Router
		zombieController   ZombieController
		zombieResolverMock *driver_zombie.IsZombieResolverMock
		configurerMock     *driver_zombie.ZombieConfigurerMock
	)

	setup := func() {
		zombieResolverMock = &driver_zombie.IsZombieResolverMock{}
		zombieController = NewZombieController(zombieResolverMock)
		configurerMock = &driver_zombie.ZombieConfigurerMock{}
		configController := NewConfigController(configurerMock)
		sut = NewRouter(zombieController, configController, NewHealthController())
	}

	t.Run("Given a server", func(t *testing.T) {
		setup()
		driverID := uuid.NewV4().String()
		zombieResolverMock.ResolveFunc = func(driverID string) bool {
			return true
		}
		configurerMock.SetZombieConfigFunc = func(in1 driver_zombie.ZombieConfigProjection) error {
			return nil
		}
		t.Run("When consuming health endpoint", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			sut.ServeHTTP(res, req)
			t.Run("Then a 200 status code is returned", func(t *testing.T) {
				assertThat.Equal(http.StatusOK, res.Code)
			})
		})
		t.Run("When consuming is zombie endpoint", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/drivers/"+driverID, nil)
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			sut.ServeHTTP(res, req)
			t.Run("Then the proper response is returned", func(t *testing.T) {
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
		t.Run("When changing the configuration", func(t *testing.T) {
			body := map[string]interface{}{"max_meters": 500, "last_minutes": 5}
			bodyJson, _ := json.Marshal(body)
			req, err := http.NewRequest(http.MethodPatch, "/config", bytes.NewBuffer(bodyJson))

			assertThat.NoError(err)
			res := httptest.NewRecorder()
			sut.ServeHTTP(res, req)
			t.Run("Then the configurer saves it", func(t *testing.T) {
				body, err := ioutil.ReadAll(res.Body)
				assertThat.NoError(err)

				assertThat.True(len(configurerMock.SetZombieConfigCalls()) > 0)

				var configResponse driver_zombie.ZombieConfigProjection
				err = json.Unmarshal(body, &configResponse)
				assertThat.NoError(err)

				assertThat.Equal(http.StatusOK, res.Code)
				assertThat.Equal(500, configResponse.MaxMeters)
				assertThat.Equal(5, configResponse.LastMinutes)
			})
		})
	})
}
