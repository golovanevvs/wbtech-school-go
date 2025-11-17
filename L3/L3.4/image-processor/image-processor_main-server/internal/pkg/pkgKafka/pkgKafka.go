package pkgKafka

import (
	"context"
	"encoding/json"

	"github.com/wb-go/wbf/kafka"
)

type ProcessImageMessage struct {
	ImageID string `json:"image_id"`
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(brokers []string, topic string) *KafkaProducer {
	producer := kafka.NewProducer(brokers, topic)
	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}
}

func (kq *KafkaProducer) SendProcessTask(ctx context.Context, imageID string) error {
	msg := ProcessImageMessage{
		ImageID: imageID,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return kq.producer.Send(ctx, nil, value)
}

func (kq *KafkaProducer) Close() error {
	return kq.producer.Close()
}
