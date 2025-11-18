package pkgKafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
	groupID  string
	rs       *retry.Strategy
}

func NewConsumer(brokers []string, topic, groupID string, rs *retry.Strategy) *KafkaConsumer {
	consumer := kafka.NewConsumer(brokers, topic, groupID)
	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
		rs:       rs,
	}
}

func (kc *KafkaConsumer) StartConsuming(ctx context.Context, handler func(msg ProcessImageMessage) error) error {
	defer log.Println("Defer: consumer stopped!!!")
	msgChan := make(chan kafkago.Message)

	kc.consumer.StartConsuming(ctx, msgChan, *kc.rs)
	log.Println("consumer started!!!")

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			var processMsg ProcessImageMessage
			err := json.Unmarshal(msg.Value, &processMsg)
			if err != nil {
				continue
			}

			err = handler(processMsg)
			if err != nil {
				continue
			}

			_ = kc.consumer.Commit(ctx, msg)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
