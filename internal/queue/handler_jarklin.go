package queue

import (
	"context"
	"microservice/internal/app/jarklin"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type HandlerJarklin struct {
	service *jarklin.Service
}

func NewHandlerJarklin(service *jarklin.Service) *HandlerJarklin {
	return &HandlerJarklin{service: service}
}

func (h *HandlerJarklin) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var found bool
	for _, header := range msg.Headers {
		zap.S().Debugf("header: %s - %s", string(header.Key), string(header.Value))
		if string(header.Key) == "type" && string(header.Value) == "NotifyResponse" {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	return h.service.HandleResponse(ctx, msg.Value)
}