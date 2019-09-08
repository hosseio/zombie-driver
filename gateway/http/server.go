package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type EmptyResponse struct{}

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

type Endpoints struct {
	RedirectEndpoints []RedirectEndpoint
	NSQEndpoints      []NSQEndpoint
}

func NewRouter(
	nsqController NSQController,
	redirectController RedirectController,
) *mux.Router {
	router := mux.NewRouter()

	for _, nsqEndpoint := range nsqController.NSQEndpoints {
		router.
			Path(nsqEndpoint.Path).
			HandlerFunc(nsqController.handleNSQ).
			Methods(nsqEndpoint.Method)
	}

	for _, redirectEndpoint := range redirectController.RedirectEndpoints {
		router.
			Path(redirectEndpoint.Path).
			HandlerFunc(redirectController.handleRedirect).
			Methods(redirectEndpoint.Method)
	}

	return router
}
