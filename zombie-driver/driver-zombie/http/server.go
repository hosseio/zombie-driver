package http

import (
	"net/http"
	"time"

	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"

	"github.com/arpando/controller"
	"github.com/gorilla/mux"
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

func NewRouter(l ZombieController) *mux.Router {
	router := mux.NewRouter()
	router.
		Path("/drivers/{id}").
		HandlerFunc(l.driverIsZombieHandler).
		Methods(http.MethodGet)

	return router
}

type ZombieController struct {
	controller.Json
	zombieResolver driver_zombie.IsZombieResolver
}

func NewLocationController(resolver driver_zombie.IsZombieResolver) ZombieController {
	return ZombieController{zombieResolver: resolver}
}

type ZombieResponse struct {
	Id     string `json:"id"`
	Zombie bool   `json:"zombie"`
}

func (l ZombieController) driverIsZombieHandler(writer http.ResponseWriter, request *http.Request) {
	l.Handle(writer, request, func() (int, interface{}) {
		vars := mux.Vars(request)
		driverID := vars["id"]

		isZombie := l.zombieResolver.Resolve(driverID)

		return http.StatusOK, ZombieResponse{driverID, isZombie}
	})
}
