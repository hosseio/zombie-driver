package bootstrap

import (
	"net/http"
	"regexp"
	"strings"

	http2 "github.com/heetch/jose-odg-technical-test/gateway/http"
)

type URLConfigToRouterConverter struct{}

func NewURLConfigToRouterConverter() URLConfigToRouterConverter {
	return URLConfigToRouterConverter{}
}

func (c URLConfigToRouterConverter) Convert(config URLConfig) http2.Endpoints {
	redirects := []http2.RedirectEndpoint{}
	nsqs := []http2.NSQEndpoint{}
	for _, urlConfig := range config.Urls {
		path := c.convertPath(urlConfig.Path)
		httpMethod := c.convertMethod(urlConfig.Method)
		if urlConfig.Http.Host != "" {
			redirectEndpoint := http2.RedirectEndpoint{
				Path:   path,
				Method: httpMethod,
				HostTo: urlConfig.Http.Host,
			}
			redirects = append(redirects, redirectEndpoint)
			continue
		}
		nsqEndpoint := http2.NSQEndpoint{
			Path:   path,
			Method: httpMethod,
			Topic:  urlConfig.Nsq.Topic,
		}
		nsqs = append(nsqs, nsqEndpoint)
	}

	return http2.Endpoints{
		NSQEndpoints:      nsqs,
		RedirectEndpoints: redirects,
	}
}

// convertPath transform params in config.yaml like :id to gorilla mux params {id}
func (c URLConfigToRouterConverter) convertPath(path string) string {
	r, _ := regexp.Compile(":[a-zA-Z]+")
	matched := r.FindAllString(path, -1)
	if matched == nil {
		return path
	}

	replaced := path
	for _, match := range matched {
		muxParam := strings.Replace(match, ":", "{", -1) + "}"
		replaced = strings.Replace(replaced, match, muxParam, -1)
	}

	return replaced
}

func (c URLConfigToRouterConverter) convertMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return http.MethodGet
	case "HEAD":
		return http.MethodHead
	case "POST":
		return http.MethodPost
	case "PUT":
		return http.MethodPut
	case "PATCH":
		return http.MethodPatch
	case "DELETE":
		return http.MethodDelete
	case "CONNECT":
		return http.MethodConnect
	case "OPTIONS":
		return http.MethodOptions
	case "TRACE":
		return http.MethodTrace
	}

	return ""
}
