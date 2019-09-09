package messaging

import (
	"encoding/json"
	"log"

	"github.com/heetch/jose-odg-technical-test/driver-location/driver-location"

	"github.com/chiguirez/cromberbus"
	"github.com/nsqio/go-nsq"
)

type CreateDriverLocationMessageDTO struct {
	DriverID string  `json:"id"`
	Lat      float32 `json:"latitude"`
	Lon      float32 `json:"longitude"`
}

type CreateDriverLocationHandler struct {
	commandBus cromberbus.CommandBus
}

func NewCreateDriverLocationHandler(commandBus cromberbus.CommandBus) CreateDriverLocationHandler {
	return CreateDriverLocationHandler{commandBus: commandBus}
}

func (c CreateDriverLocationHandler) HandleMessage(message *nsq.Message) error {
	createDriverLocationMessage := CreateDriverLocationMessageDTO{}
	err := json.Unmarshal(message.Body, &createDriverLocationMessage)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	command := driver_location.CreateLocationCommand{
		Latitude:  float64(createDriverLocationMessage.Lat),
		Longitude: float64(createDriverLocationMessage.Lon),
		DriverID:  createDriverLocationMessage.DriverID,
	}

	return c.commandBus.Dispatch(command)
}
