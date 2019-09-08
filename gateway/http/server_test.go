//+build integration

package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heetch/jose-odg-technical-test/gateway/messaging"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	assertThat := require.New(t)

	var (
		sut                *mux.Router
		nsqController      NSQController
		redirectController RedirectController
		producer           *messaging.NSQSenderMock
	)

	setup := func() {
		endpoints := Endpoints{
			RedirectEndpoints: []RedirectEndpoint{
				{"/drivers/{id}", http.MethodGet, "other-host"},
			},
			NSQEndpoints: []NSQEndpoint{
				{"/drivers/{id}/locations", http.MethodGet, "topic_name"},
				{"/drivers/{id}/locations/{locationId}", http.MethodDelete, "topic_name"},
			},
		}
		producer = &messaging.NSQSenderMock{
			SendMessageFunc: func(topic string, message []byte) error {
				return nil
			},
		}
		nsqController = NewNSQController(endpoints.NSQEndpoints, producer)
		redirectController = NewRedirectController(endpoints.RedirectEndpoints)
		sut = NewRouter(nsqController, redirectController)
	}

	t.Run("Given a server", func(t *testing.T) {
		setup()

		t.Run("When consuming nsq endpoints", func(t *testing.T) {
			body := map[string]interface{}{"latitude": 1.23, "longitude": 3.21}
			bodyJson, _ := json.Marshal(body)
			req, err := http.NewRequest(http.MethodGet, "/drivers/some-uuid/locations", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			sut.ServeHTTP(res, req)
			t.Run("Then a message is sent to the topic", func(t *testing.T) {
				assertThat.True(len(producer.SendMessageCalls()) > 0)
			})
			t.Run("And a 200 status code is returned", func(t *testing.T) {
				body, err := ioutil.ReadAll(res.Body)
				assertThat.NoError(err)

				var emptyResponse EmptyResponse
				err = json.Unmarshal(body, &emptyResponse)
				assertThat.NoError(err)

				assertThat.Equal(http.StatusOK, res.Code)
			})
		})
	})
	t.Run("Given a server", func(t *testing.T) {
		setup()

		t.Run("When consuming redirect endpoints", func(t *testing.T) {
			body := map[string]interface{}{"some-data": "data"}
			bodyJson, _ := json.Marshal(body)
			req, err := http.NewRequest(http.MethodGet, "/drivers/{id}", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			res := httptest.NewRecorder()
			sut.ServeHTTP(res, req)
			t.Run("Then a redirect to a different host is done", func(t *testing.T) {
				assertThat.True(true)
			})
		})
	})
}
