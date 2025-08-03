package kafka

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Config struct {
	Brokers             []string
	ClientID            string
	ConsumerGroup       string
	Version             string
	RetryMax            int
	RequiredAcks        sarama.RequiredAcks
	Partitioner         string
	EnableReturnSuccess bool
}

func DefaultConfig(brokers []string, group string) *Config {
	return &Config{
		Brokers:             brokers,
		ClientID:            "go-kafka-client",
		ConsumerGroup:       group,
		Version:             "2.8.0",
		RetryMax:            5,
		RequiredAcks:        sarama.WaitForAll,
		Partitioner:         "hash",
		EnableReturnSuccess: true,
	}
}

func (c *Config) buildSaramaConfig() (*sarama.Config, error) {
	saramaCfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, err
	}
	saramaCfg.Version = version

	saramaCfg.ClientID = c.ClientID

	saramaCfg.Producer.Retry.Max = c.RetryMax
	saramaCfg.Producer.RequiredAcks = c.RequiredAcks
	saramaCfg.Producer.Return.Successes = c.EnableReturnSuccess
	saramaCfg.Producer.Return.Errors = true

	switch c.Partitioner {
	case "roundrobin":
		saramaCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "hash":
		saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
	default:
		saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
	}

	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	return saramaCfg, nil
}

type Producer struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	config        *Config
}

func NewSyncProducer(cfg *Config) (*Producer, error) {
	saramaCfg, err := cfg.buildSaramaConfig()
	if err != nil {
		return nil, err
	}

	syncProd, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}
	return &Producer{
		syncProducer: syncProd,
		config:       cfg,
	}, nil
}

func NewAsyncProducer(cfg *Config) (*Producer, error) {
	saramaCfg, err := cfg.buildSaramaConfig()
	if err != nil {
		return nil, err
	}
	asyncProd, err := sarama.NewAsyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}
	return &Producer{
		asyncProducer: asyncProd,
		config:        cfg,
	}, nil
}

func (p *Producer) SendSync(ctx context.Context, topic string, key, value []byte, headers map[string][]byte) (partition int32, offset int64, err error) {
	if p.syncProducer == nil {
		return 0, 0, errors.New("sync producer не инициализирован")
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	if headers != nil {
		for k, v := range headers {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: v,
			})
		}
	}
	return p.syncProducer.SendMessage(msg)
}

func (p *Producer) SendAsync(topic string, key, value []byte, headers map[string][]byte) error {
	if p.asyncProducer == nil {
		return errors.New("async producer не инициализирован")
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	if headers != nil {
		for k, v := range headers {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: v,
			})
		}
	}
	p.asyncProducer.Input() <- msg
	return nil
}

func (p *Producer) Close() error {
	var err error
	if p.syncProducer != nil {
		if e := p.syncProducer.Close(); e != nil {
			err = e
		}
	}
	if p.asyncProducer != nil {
		if e := p.asyncProducer.Close(); e != nil {
			err = e
		}
	}
	return err
}

type MessageHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type consumerGroupHandler struct {
	handler MessageHandler
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim — основной цикл получения сообщений
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := context.Background()
		if err := c.handler(ctx, msg); err != nil {
			// логика: можно логировать и/или пропускать
			log.Printf("ошибка при обработке сообщения: %v", err)
			// не подтверждаем, чтобы оно могло переобработаться в зависимости от стратегии
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

type ConsumerGroup struct {
	group   sarama.ConsumerGroup
	config  *Config
	handler MessageHandler
	topics  []string
}

func NewConsumerGroup(cfg *Config, topics []string, handler MessageHandler) (*ConsumerGroup, error) {
	saramaCfg, err := cfg.buildSaramaConfig()
	if err != nil {
		return nil, err
	}
	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, err
	}
	return &ConsumerGroup{
		group:   group,
		config:  cfg,
		handler: handler,
		topics:  topics,
	}, nil
}

func (c *ConsumerGroup) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	errCh := make(chan error, 1)

	go func() {
		defer wg.Done()
		for {
			if ctx.Err() != nil {
				return
			}
			handler := &consumerGroupHandler{handler: c.handler}
			if err := c.group.Consume(ctx, c.topics, handler); err != nil {
				log.Printf("ошибка consumer group consume: %v", err)
				// если фатально — передать вверх
				errCh <- err
				return
			}
			// после успешного Consume, проверить, не отменён ли контекст и повторить (ребалансировки)
		}
	}()

	// Ждём или ошибку
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// ждём горутину завершить аккуратно
		wg.Wait()
		return nil
	}
}

// Close закрывает consumer group
func (c *ConsumerGroup) Close() error {
	return c.group.Close()
}
