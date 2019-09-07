//+build integration

package messaging

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/nsqio/go-nsq"

	"github.com/stretchr/testify/require"
)

func TestNsqConsumer(t *testing.T) {
	assertThat := require.New(t)
	nsqAddress := "127.0.0.1:4150"
	topic := uuid.NewV4().String()

	var (
		handler     nsq.HandlerFunc
		messageChan chan string
		producer    *nsq.Producer
	)

	setup := func() {
		messageChan = make(chan string, 0)

		config := nsq.NewConfig()
		producer, _ = nsq.NewProducer(nsqAddress, config)

		handler = nsq.HandlerFunc(func(message *nsq.Message) error {
			messageChan <- string(message.Body)
			return nil
		})
	}
	t.Run("Given a NsqConsumer running", func(t *testing.T) {
		setup()
		defer producer.Stop()

		sut := NewNsqConsumer(NsqAddr(nsqAddress), TopicAddr(topic), ChannelAddr("ch"), handler)
		go func() {
			err := sut.Run(context.Background())
			assertThat.NoError(err)
		}()
		t.Run("When a message is sent to that topic", func(t *testing.T) {
			err := producer.Publish(topic, []byte("this is the message"))
			assertThat.NoError(err)
			t.Run("Then the message is read by the consumer", func(t *testing.T) {
				read := <-messageChan
				assertThat.Equal("this is the message", read)
			})
		})
	})
}
