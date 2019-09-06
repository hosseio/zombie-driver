package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/satori/go.uuid"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	assertThat := require.New(t)

	var (
		//sut                            *http.Server
		router                         *mux.Router
		locationController             LocationController
		locationsByDriverAndTimeGetter *driver_location.LocationsByDriverAndTimeGetterMock
	)

	setup := func() {
		locationsByDriverAndTimeGetter = &driver_location.LocationsByDriverAndTimeGetterMock{}
		locationController = NewLocationController(locationsByDriverAndTimeGetter)
		router = NewRouter(locationController)
		//sut = NewServer(ServerAddr(":0"), router)
	}

	t.Run("Given a server", func(t *testing.T) {
		setup()
		//go sut.ListenAndServe()
		updatedAt, _ := time.Parse(time.RFC3339, "2018-04-05T22:36:16Z")
		secondUpdatedAt, _ := time.Parse(time.RFC3339, "2018-04-05T22:36:21Z")
		driverID := uuid.NewV4().String()
		locationsByDriverAndTimeGetter.GetFunc = func(driverID string, from time.Time) ([]driver_location.LocationDTO, error) {
			return []driver_location.LocationDTO{
				{48.864193, 2.350498, updatedAt},
				{48.863921, 2.349211, secondUpdatedAt},
			}, nil
		}
		t.Run("When consuming driver locations endpoint", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/drivers/"+driverID+"/locations?minutes=5", nil)
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			t.Run("Then the locations are retrieved", func(t *testing.T) {
				body, err := ioutil.ReadAll(res.Body)
				assertThat.NoError(err)

				var locations []driver_location.LocationDTO
				err = json.Unmarshal(body, &locations)
				assertThat.NoError(err)

				assertThat.Equal(http.StatusOK, res.Code)
				assertThat.Equal(2, len(locations))
			})
		})
	})
}
