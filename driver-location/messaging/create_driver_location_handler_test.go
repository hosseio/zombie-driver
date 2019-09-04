//+build integration

package messaging

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/chiguirez/cromberbus"

	"github.com/heetch/jose-odg-technical-test/pkg"
	"github.com/satori/go.uuid"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/require"
)

func TestCreateDriverLocationConsumer(t *testing.T) {
	assertThat := require.New(t)
	nsqAddress := "127.0.0.1:4150"
	topic := "topic_name"

	var (
		producer *nsq.Producer
	)

	setup := func() {
		config := nsq.NewConfig()
		producer, _ = nsq.NewProducer(nsqAddress, config)
	}

	tearDown := func() {
		producer.Stop()
		nsq.UnRegister(topic, "")
	}
	t.Run("Given a NsqConsumer running with the create driver location consumer", func(t *testing.T) {
		setup()
		defer tearDown()

		bus := &CommandBusMock{
			DispatchFunc: func(command cromberbus.Command) error {
				return nil
			},
		}

		sut := NewCreateDriverLocationHandler(bus)
		consumer := NewNsqConsumer(nsqAddress, topic, "ch", sut)
		go func() {
			err := consumer.Run(context.Background())
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
				assertThat.True(len(bus.DispatchCalls()) > 0)
			})
		})
	})
}
