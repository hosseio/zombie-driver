package http

import (
	"encoding/json"
	"io/ioutil"
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
			return http.StatusInternalServerError, err.Error()
		}
		hostTo := c.getHostByPath(path)

		url := *request.URL
		url.Host = hostTo
		url.Scheme = "http"

		newRequest, err := http.NewRequest(request.Method, url.String(), request.Body)
		newRequest.URL = &url
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		response, err := http.DefaultClient.Do(newRequest)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		jsonResponse, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		b := map[string]interface{}{}
		err = json.Unmarshal(jsonResponse, &b)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		return response.StatusCode, b
	})
}
