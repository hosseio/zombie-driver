package http

import (
	"net/http"

	"github.com/arpando/controller"
)

type HealthController struct {
	controller.Json
}

func NewHealthController() HealthController {
	return HealthController{}
}

func (c HealthController) healthz(w http.ResponseWriter, r *http.Request) {
	c.Handle(w, r, func() (int, interface{}) {
		return http.StatusOK, "OK"
	})
}
