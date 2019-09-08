// +build wireinject

package bootstrap

import (
	netHttp "net/http"

	"github.com/heetch/jose-odg-technical-test/gateway/http"
	"github.com/heetch/jose-odg-technical-test/gateway/messaging"

	"github.com/google/wire"
)

var HttpSet = wire.NewSet(
	http.NewServer,
	http.NewRouter,
	http.NewNSQController,
	http.NewRedirectController,
)

var MessagingSet = wire.NewSet(
	messaging.NewNSQPRoducer,
)

func InitializeServer(cfg Config) (*netHttp.Server, error) {
	wire.Build(
		HttpSet,
		serverAddr,
		InitializeProducer,
		getNSQEndpoints,
		getRedirectEndpoints,
		wire.Bind(new(messaging.NSQSender), messaging.NSQProducer{}),
	)

	return &netHttp.Server{}, nil
}

func getNSQEndpoints(cfg Config) []http.NSQEndpoint {
	endpoints := getEndpoints(cfg)

	return endpoints.NSQEndpoints
}

func getRedirectEndpoints(cfg Config) []http.RedirectEndpoint {
	endpoints := getEndpoints(cfg)

	return endpoints.RedirectEndpoints
}

var endpoints *http.Endpoints

func getEndpoints(cfg Config) http.Endpoints {
	if endpoints != nil {
		return *endpoints
	}
	converter := NewURLConfigToRouterConverter()

	convert := converter.Convert(cfg.UrlConfig())
	endpoints = &convert

	return *endpoints
}

func serverAddr(cfg Config) http.ServerAddr {
	return http.ServerAddr(cfg.Server.Addr)
}

func InitializeProducer(cfg Config) (messaging.NSQProducer, error) {
	wire.Build(
		MessagingSet,
		getNsqAddr,
	)

	return messaging.NSQProducer{}, nil
}

func getNsqAddr(cfg Config) messaging.NsqAddr {
	return messaging.NsqAddr(cfg.Nsq.Addr)
}
