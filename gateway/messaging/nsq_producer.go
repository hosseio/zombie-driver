package messaging

import (
	"github.com/nsqio/go-nsq"
)

//go:generate moq -out nsq_sender_mock.go . NSQSender
type NSQSender interface {
	SendMessage(topic string, message []byte) error
}

type NSQProducer struct {
	producer *nsq.Producer
}

type NsqAddr string

func NewNSQPRoducer(address NsqAddr) (NSQProducer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(string(address), config)
	if err != nil {
		return NSQProducer{}, err
	}

	return NSQProducer{
		producer: producer,
	}, nil
}

func (p NSQProducer) SendMessage(topic string, message []byte) error {
	return p.producer.Publish(topic, message)
}
