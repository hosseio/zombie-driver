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

func NewNsqConsumer(address, topic, channel string, handler nsq.Handler) NsqConsumer {
	config := nsq.NewConfig()
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	consumer.AddHandler(handler)

	return NsqConsumer{
		consumer: consumer,
		addr:     address,
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

	ctxCancel()
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
