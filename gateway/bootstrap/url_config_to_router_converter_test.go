package bootstrap

import (
	"net/http"
	"testing"

	"github.com/heetch/jose-odg-technical-test/gateway/bootstrap"

	"github.com/stretchr/testify/require"
)

func TestURLConfigToRouterConverter_Convert(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given an URL config struct", func(t *testing.T) {
		urlConfig := bootstrap.URLConfig{
			Urls: []bootstrap.Url{
				bootstrap.Url{
					Path:   "/drivers/:id",
					Method: "GET",
					Http:   bootstrap.Http{Host: "other-host"},
				},
				bootstrap.Url{
					Path:   "/drivers/:id/locations",
					Method: "GET",
					Nsq:    bootstrap.Nsq{Topic: "topic_name"},
				},
				bootstrap.Url{
					Path:   "/drivers/:id/locations/:locationId",
					Method: "DELETE",
					Nsq:    bootstrap.Nsq{Topic: "topic_name"},
				},
			},
		}
		sut := NewURLConfigToRouterConverter()
		t.Run("When converting them to router endpoints", func(t *testing.T) {
			endpoints := sut.Convert(urlConfig)
			t.Run("Then the conversion is done properly", func(t *testing.T) {
				expectedEndpoints := Endpoints{
					RedirectEndpoints: []RedirectEndpoint{
						{"/drivers/{id}", http.MethodGet, "other-host"},
					},
					NSQEndpoints: []NSQEndpoint{
						{"/drivers/{id}/locations", http.MethodGet, "topic_name"},
						{"/drivers/{id}/locations/{locationId}", http.MethodDelete, "topic_name"},
					},
				}
				assertThat.Equal(expectedEndpoints, endpoints)
			})
		})
	})
}
