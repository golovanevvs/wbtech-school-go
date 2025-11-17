package pkgKafka

import (
	"context"
	"encoding/json"

	"github.com/wb-go/wbf/kafka"
)

type ProcessImageMessage struct {
	ImageID string `json:"image_id"`
}

type Kafka struct {
	producer *kafka.Producer
	topic    string
}

func New(brokers []string, topic string) *Kafka {
	producer := kafka.NewProducer(brokers, topic)
	return &Kafka{
		producer: producer,
		topic:    topic,
	}
}

func (kq *Kafka) SendProcessTask(ctx context.Context, imageID string) error {
	msg := ProcessImageMessage{
		ImageID: imageID,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return kq.producer.Send(ctx, nil, value)
}

func (kq *Kafka) Close() error {
	return kq.producer.Close()
}
