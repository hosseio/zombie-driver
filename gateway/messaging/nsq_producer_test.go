//+build integration

package messaging

import (
	"os"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNsqConsumer(t *testing.T) {
	assertThat := require.New(t)
	nsqAddress := os.Getenv("NSQ")
	if nsqAddress == "" {
		nsqAddress = "127.0.0.1:4150"
	}
	topic := uuid.NewV4().String()

	var (
		handler     nsq.HandlerFunc
		messageChan chan string
		sut         NSQProducer
		consumer    *nsq.Consumer
	)

	setup := func() {
		messageChan = make(chan string, 0)

		sut, _ = NewNSQPRoducer(NsqAddr(nsqAddress))

		config := nsq.NewConfig()
		consumer, _ = nsq.NewConsumer(topic, "ch", config)
		handler = nsq.HandlerFunc(func(message *nsq.Message) error {
			messageChan <- string(message.Body)
			return nil
		})
		consumer.AddHandler(handler)
		consumer.ConnectToNSQD(nsqAddress)
	}
	t.Run("Given a NSQProducer", func(t *testing.T) {
		setup()
		t.Run("When a message is sent to the topic", func(t *testing.T) {
			err := sut.SendMessage(topic, []byte("this is the message"))
			assertThat.NoError(err)
			t.Run("Then the message can be read", func(t *testing.T) {
				read := <-messageChan
				assertThat.Equal("this is the message", read)
			})
		})
	})
}
