package queue

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type Consumer struct {
	ready chan bool
	finish chan bool
	cancel context.CancelFunc

	client sarama.ConsumerGroup
	handler Handler
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				zap.S().Infow("message channel was closed")

				return nil
			}
			if err := c.handler.Handle(session.Context(), message); err != nil {
				zap.S().Errorw("consumer handle message", "error", err, "topic", message.Topic, "partition", message.Partition, "offset", message.Offset)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}


func NewConsumer(cxt context.Context, brokers []string, groupName string, handler Handler) (*Consumer, error) {
	version, err := sarama.ParseKafkaVersion("3.2.2")
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, groupName, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		ready: make(chan bool),
		client: client,
		handler: handler,
	}, nil
}

func (c *Consumer) Run(ctx context.Context, topics []string) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	go func() {
		for {
			if err := c.client.Consume(ctx, topics, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				zap.S().Fatalw("consume failed", "error", err)
			}
			if ctx.Err() != nil {
				close(c.finish)

				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready
}

func (c *Consumer) Stop() {
	c.cancel()
	<-c.finish
}
