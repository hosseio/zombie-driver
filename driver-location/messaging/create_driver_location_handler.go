package messaging

import (
	"encoding/json"

	"github.com/chiguirez/cromberbus"
	"github.com/heetch/jose-odg-technical-test/driver-location"
	"github.com/heetch/jose-odg-technical-test/pkg"
	"github.com/nsqio/go-nsq"
)

type CreateDriverLocationHandler struct {
	commandBus cromberbus.CommandBus
}

func NewCreateDriverLocationHandler(commandBus cromberbus.CommandBus) CreateDriverLocationHandler {
	return CreateDriverLocationHandler{commandBus: commandBus}
}

func (c CreateDriverLocationHandler) HandleMessage(message *nsq.Message) error {
	createDriverLocationMessage := pkg.CreateDriverLocationMessage{}
	err := json.Unmarshal(message.Body, &createDriverLocationMessage)
	if err != nil {
		return err
	}

	command := driver_location.CreateLocationCommand{
		Latitude:  float64(createDriverLocationMessage.Lat),
		Longitude: float64(createDriverLocationMessage.Lon),
		DriverID:  createDriverLocationMessage.DriverID,
	}

	return c.commandBus.Dispatch(command)
}
