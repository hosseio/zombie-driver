package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/heetch/jose-odg-technical-test/gateway/messaging"

	"github.com/gorilla/mux"

	"github.com/arpando/controller"
)

type NSQEndpoint struct {
	Path   string
	Method string
	Topic  string
}

type NSQController struct {
	controller.Json
	NSQEndpoints []NSQEndpoint
	producer     messaging.NSQSender
}

func NewNSQController(endpoints []NSQEndpoint, producer messaging.NSQSender) NSQController {
	return NSQController{
		NSQEndpoints: endpoints,
		producer:     producer,
	}
}

func (c NSQController) getTopicByPath(path string) string {
	for _, endpoint := range c.NSQEndpoints {
		if endpoint.Path == path {
			return endpoint.Topic
		}
	}

	return ""
}

func (c NSQController) handleNSQ(writer http.ResponseWriter, request *http.Request) {
	c.Handle(writer, request, func() (int, interface{}) {
		path, err := mux.CurrentRoute(request).GetPathTemplate()
		if err != nil {
			return http.StatusInternalServerError, EmptyResponse{}
		}
		topic := c.getTopicByPath(path)

		message, err := c.message(request)
		if err != nil {
			return http.StatusInternalServerError, EmptyResponse{}
		}

		err = c.producer.SendMessage(topic, message)
		if err != nil {
			return http.StatusInternalServerError, EmptyResponse{}
		}

		return http.StatusOK, EmptyResponse{}
	})
}

func (c NSQController) message(request *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(request.Body)
	var b map[string]interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		return body, err
	}

	vars := mux.Vars(request)
	for i, v := range vars {
		b[i] = v
	}

	return json.Marshal(b)
}
