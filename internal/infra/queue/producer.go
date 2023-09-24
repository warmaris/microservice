package queue

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	conn sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	pr, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{conn: pr}, nil
}

func (p *Producer) Produce(topic, key string, value []byte, msgType string) error {
	partition, offset, err := p.conn.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{
				Key: []byte("type"),
				Value: []byte(msgType),
			},
		},
	})

	if err == nil {
		zap.S().Debugw("kafka send message", "partition", partition, "offset", offset)
	} else {
		zap.S().Errorw("kafka send error", "error", err)
	}

	return err
}
