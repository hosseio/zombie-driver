package messaging

import (
	"context"
	"log"

	"github.com/nsqio/go-nsq"
	"golang.org/x/sync/errgroup"
)

type NsqConsumer struct {
	consumer *nsq.Consumer
	addr     string
}

type NsqAddr string
type TopicAddr string
type ChannelAddr string

func NewNsqConsumer(address NsqAddr, topic TopicAddr, channel ChannelAddr, handler nsq.Handler) NsqConsumer {
	config := nsq.NewConfig()
	consumer, _ := nsq.NewConsumer(string(topic), string(channel), config)
	consumer.AddHandler(handler)

	return NsqConsumer{
		consumer: consumer,
		addr:     string(address),
	}
}

func (c NsqConsumer) Run(ctx context.Context) error {
	ctx, ctxCancel := context.WithCancel(ctx)

	g, ctx := errgroup.WithContext(ctx)
	err := c.consumer.ConnectToNSQD(c.addr)
	if err != nil {
		log.Panic("Could not connect")
	}
	g.Go(func() error {
		for {
			<-c.consumer.StopChan
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
	})

	<-ctx.Done()
	ctxCancel()
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
