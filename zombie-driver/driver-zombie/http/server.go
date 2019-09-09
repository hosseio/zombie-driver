package http

import (
	"encoding/json"
	"io/ioutil"
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

func NewRouter(zc ZombieController, cc ConfigController, hc HealthController) *mux.Router {
	router := mux.NewRouter()
	router.
		Path("/drivers/{id}").
		HandlerFunc(zc.driverIsZombieHandler).
		Methods(http.MethodGet)
	router.Path("/config").
		HandlerFunc(cc.updateConfigHandler).
		Methods(http.MethodPatch)

	router.Path("/healthz").Methods(http.MethodGet).HandlerFunc(hc.healthz)

	return router
}

type ConfigController struct {
	controller.Json
	configurer driver_zombie.ZombieConfigurer
}

func NewConfigController(configurer driver_zombie.ZombieConfigurer) ConfigController {
	return ConfigController{configurer: configurer}
}

func (c ConfigController) updateConfigHandler(writer http.ResponseWriter, request *http.Request) {
	c.Handle(writer, request, func() (int, interface{}) {
		body, err := ioutil.ReadAll(request.Body)
		var projection driver_zombie.ZombieConfigProjection
		err = json.Unmarshal(body, &projection)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = c.configurer.SetZombieConfig(projection)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, projection
	})
}

type ZombieController struct {
	controller.Json
	zombieResolver driver_zombie.IsZombieResolver
}

func NewZombieController(resolver driver_zombie.IsZombieResolver) ZombieController {
	return ZombieController{zombieResolver: resolver}
}

type ZombieResponse struct {
	Id     string `json:"id"`
	Zombie bool   `json:"zombie"`
}

func (c ZombieController) driverIsZombieHandler(writer http.ResponseWriter, request *http.Request) {
	c.Handle(writer, request, func() (int, interface{}) {
		vars := mux.Vars(request)
		driverID := vars["id"]

		isZombie := c.zombieResolver.Resolve(driverID)

		return http.StatusOK, ZombieResponse{driverID, isZombie}
	})
}
