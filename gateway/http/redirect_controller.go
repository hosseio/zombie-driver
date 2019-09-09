package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arpando/controller"
)

type RedirectEndpoint struct {
	Path   string
	Method string
	HostTo string
}

type RedirectController struct {
	controller.Json
	RedirectEndpoints []RedirectEndpoint
}

func NewRedirectController(endpoints []RedirectEndpoint) RedirectController {
	return RedirectController{
		RedirectEndpoints: endpoints,
	}
}

func (c RedirectController) getHostByPath(path string) string {
	for _, endpoint := range c.RedirectEndpoints {
		if endpoint.Path == path {
			return endpoint.HostTo
		}
	}

	return ""
}

func (c RedirectController) handleRedirect(writer http.ResponseWriter, request *http.Request) {
	c.Handle(writer, request, func() (int, interface{}) {
		route := mux.CurrentRoute(request)
		path, err := route.GetPathTemplate()
		if err != nil {
			return http.StatusInternalServerError, EmptyResponse{}
		}
		hostTo := c.getHostByPath(path)

		http.Redirect(writer, request, hostTo, http.StatusMovedPermanently)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return http.StatusInternalServerError, EmptyResponse{}
		}

		return response.StatusCode, response.Body
	})
}
