package jarklin

import (
	"context"
	"encoding/json"
	"microservice/internal/infra/queue"
)

const TopicJarklinEvents = "jarklin_events"

type kafkaSender struct {
	producer *queue.Producer
	topic string
}

func NewSender(producer *queue.Producer) Sender {
	return &kafkaSender{producer: producer, topic: TopicJarklinEvents}
}

func (k *kafkaSender) send(ctx context.Context, ev NotifyRequest) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	return k.producer.Produce(k.topic, ev.MessageID, payload, "NotifyRequest")
}
