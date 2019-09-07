//+build integration

package messaging

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/chiguirez/cromberbus"

	"github.com/heetch/jose-odg-technical-test/driver-location/pkg"
	"github.com/satori/go.uuid"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/require"
)

func TestCreateDriverLocationConsumer(t *testing.T) {
	assertThat := require.New(t)
	nsqAddress := "127.0.0.1:4150"
	topic := uuid.NewV4().String()

	var (
		sut CreateDriverLocationHandler

		consumer NsqConsumer
		bus      *CommandBusMock
		producer *nsq.Producer
		syncChan chan bool
	)

	setup := func() {
		config := nsq.NewConfig()
		producer, _ = nsq.NewProducer(nsqAddress, config)

		syncChan = make(chan bool, 0)
		bus = &CommandBusMock{
			DispatchFunc: func(command cromberbus.Command) error {
				syncChan <- true
				return nil
			},
		}

		sut = NewCreateDriverLocationHandler(bus)
		consumer = NewNsqConsumer(NsqAddr(nsqAddress), TopicAddr(topic), ChannelAddr("ch"), sut)

		syncChan = make(chan bool, 0)
	}

	tearDown := func() {
		producer.Stop()
		nsq.UnRegister(topic, "")
	}
	t.Run("Given a NsqConsumer running with the create driver location consumer", func(t *testing.T) {
		setup()
		defer tearDown()

		ctx := context.Background()
		go func() {
			err := consumer.Run(ctx)
			assertThat.NoError(err)
		}()
		t.Run("When a create driver location message is sent to that topic", func(t *testing.T) {
			message := pkg.CreateDriverLocationMessage{
				DriverID: uuid.NewV4().String(),
				Lat:      0.0,
				Lon:      0.0,
			}
			bytes, err := json.Marshal(message)
			assertThat.NoError(err)

			err = producer.Publish(topic, bytes)
			assertThat.NoError(err)
			t.Run("Then the command bus dispatch a command", func(t *testing.T) {
				<-syncChan
				assertThat.True(len(bus.DispatchCalls()) > 0)
			})
		})
	})
}
