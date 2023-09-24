package stubs

import (
	"context"
	"encoding/json"
	"math/rand"
	"microservice/internal/app/jarklin"
	"microservice/internal/infra/queue"

	"github.com/IBM/sarama"
)

type JarklinConsumerHandler struct {
	producer *queue.Producer
}

func NewJarklinConsumerHandler(producer *queue.Producer) *JarklinConsumerHandler {
	return &JarklinConsumerHandler{producer: producer}
}

func (h *JarklinConsumerHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var found bool
	for _, header := range msg.Headers {
		if string(header.Key) == "type" && string(header.Value) == "NotifyRequest" {
			found = true
			break
		}
	}
	if !found {
		return nil
	}
	
	req := new(jarklin.NotifyRequest)
	err := json.Unmarshal(msg.Value, req)
	if err != nil {
		return err
	}

	res := &jarklin.NotifyResponse{
		MessageID: req.MessageID,
		Status: jarklin.StatusSuccess,
	}

	if rand.Intn(100) >= 90 {
		res.Status = jarklin.StatusFail
		res.StatusInfo = "stub internal error"
	}

	payload, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return h.producer.Produce(jarklin.TopicJarklinEvents, res.MessageID, payload, "NotifyResponse")
}