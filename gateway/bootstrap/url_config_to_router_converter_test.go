package bootstrap

import (
	"net/http"
	"testing"

	http2 "github.com/heetch/jose-odg-technical-test/gateway/http"

	"github.com/stretchr/testify/require"
)

func TestURLConfigToRouterConverter_Convert(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given an URL config struct", func(t *testing.T) {
		urlConfig := URLConfig{
			Urls: []Url{
				Url{
					Path:   "/drivers/:id",
					Method: "GET",
					Http:   Http{Host: "other-host"},
				},
				Url{
					Path:   "/drivers/:id/locations",
					Method: "GET",
					Nsq:    Nsq{Topic: "topic_name"},
				},
				Url{
					Path:   "/drivers/:id/locations/:locationId",
					Method: "DELETE",
					Nsq:    Nsq{Topic: "topic_name"},
				},
			},
		}
		sut := NewURLConfigToRouterConverter()
		t.Run("When converting them to router endpoints", func(t *testing.T) {
			endpoints := sut.Convert(urlConfig)
			t.Run("Then the conversion is done properly", func(t *testing.T) {
				expectedEndpoints := http2.Endpoints{
					RedirectEndpoints: []http2.RedirectEndpoint{
						{"/drivers/{id}", http.MethodGet, "other-host"},
					},
					NSQEndpoints: []http2.NSQEndpoint{
						{"/drivers/{id}/locations", http.MethodGet, "topic_name"},
						{"/drivers/{id}/locations/{locationId}", http.MethodDelete, "topic_name"},
					},
				}
				assertThat.Equal(expectedEndpoints, endpoints)
			})
		})
	})
}
