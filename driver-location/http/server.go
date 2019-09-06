package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arpando/controller"
	"github.com/gorilla/mux"
	"github.com/heetch/jose-odg-technical-test/driver-location"
)

type ServerAddr string

func NewServer(addr ServerAddr, router *mux.Router) *http.Server {
	return &http.Server{
		Addr:         string(addr),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: router,
	}
}

func NewRouter(l LocationController) *mux.Router {
	router := mux.NewRouter()
	router.
		Path("/drivers/{id}/locations").
		HandlerFunc(l.getLocationsHandler).
		Methods(http.MethodGet)

	return router
}

type LocationController struct {
	controller.Json
	queryService driver_location.LocationsByDriverAndTimeGetter
}

func NewLocationController(queryService driver_location.LocationsByDriverAndTimeGetter) LocationController {
	return LocationController{queryService: queryService}
}

func (l LocationController) getLocationsHandler(writer http.ResponseWriter, request *http.Request) {
	l.Handle(writer, request, func() (int, interface{}) {
		vars := mux.Vars(request)
		driverID := vars["id"]
		minutes, _ := strconv.Atoi(request.FormValue("minutes"))
		now := time.Now()
		from := now.Add(time.Duration(-minutes) * time.Minute)

		locations, err := l.queryService.Get(driverID, from)
		if err != nil {
			return http.StatusNotFound, locations
		}

		return http.StatusOK, locations
	})
}
