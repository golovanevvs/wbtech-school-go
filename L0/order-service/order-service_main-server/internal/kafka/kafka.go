package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/IBM/sarama"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/config"
	"github.com/rs/zerolog"
)

type Kafka struct {
	cfg                 *sarama.Config
	brokers             []string
	consumerGroup       string
	consumerWorkerCount int
	log                 *zerolog.Logger
}

func New(cfg *config.Kafka, logger *zerolog.Logger) (*Kafka, error) {
	log := logger.With().Str("component", "kafka").Logger()

	saramaCfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse Kafka version")
		return nil, err
	}
	saramaCfg.Version = version
	saramaCfg.ClientID = cfg.ClientID
	saramaCfg.Producer.Retry.Max = cfg.RetryMax
	saramaCfg.Producer.RequiredAcks = sarama.RequiredAcks(cfg.RequiredAcks)
	saramaCfg.Producer.Return.Successes = cfg.EnableReturnSuccess
	saramaCfg.Producer.Return.Errors = true

	switch cfg.Partitioner {
	case "roundrobin":
		saramaCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "hash":
		saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
	default:
		log.Warn().Str("partitioner", cfg.Partitioner).Msg("unknown partitioner, defaulting to hash")
		saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
	}

	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	return &Kafka{
		cfg:                 saramaCfg,
		brokers:             cfg.Brokers,
		consumerGroup:       cfg.ConsumerGroup,
		consumerWorkerCount: cfg.ConsumerWorkerCount,
		log:                 &log,
	}, nil
}

type Producer struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	log           *zerolog.Logger
	closeWg       sync.WaitGroup
}

func NewSyncProducer(k *Kafka) (*Producer, error) {
	log := k.log.With().Str("component", "sync producer").Logger()

	syncProd, err := sarama.NewSyncProducer(k.brokers, k.cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create sync producer")
		return nil, err
	}

	return &Producer{
		syncProducer: syncProd,
		log:          &log,
	}, nil
}

func NewAsyncProducer(k *Kafka) (*Producer, error) {
	log := k.log.With().Str("component", "async producer").Logger()

	asyncProd, err := sarama.NewAsyncProducer(k.brokers, k.cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create async producer")
		return nil, err
	}

	p := &Producer{
		asyncProducer: asyncProd,
		log:           &log,
	}

	p.closeWg.Add(2)

	go func() {
		defer p.closeWg.Done()
		for err := range asyncProd.Errors() {
			log.Error().Err(err.Err).
				Str("topic", err.Msg.Topic).
				Msg("async send error")
		}
	}()

	go func() {
		defer p.closeWg.Done()
		for msg := range asyncProd.Successes() {
			log.Debug().
				Str("topic", msg.Topic).
				Msg("async message sent successfully")
		}
	}()

	return p, nil
}

func (p *Producer) buildMessage(topic string, key, value []byte, headers map[string][]byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	for k, v := range headers {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: v,
		})
	}
	return msg
}

func (p *Producer) SendSync(ctx context.Context, topic string, key, value []byte, headers map[string][]byte) (partition int32, offset int64, err error) {
	if p.syncProducer == nil {
		p.log.Error().Err(err).
			Str("topic", topic).
			Msg("sync producer is not initialized")
		return 0, 0, errors.New("sync producer is not initialized")
	}

	msg := p.buildMessage(topic, key, value, headers)

	resultCh := make(chan struct{})
	var part int32
	var off int64
	var sendErr error

	go func() {
		defer close(resultCh)
		part, off, sendErr = p.syncProducer.SendMessage(msg)
	}()

	select {
	case <-ctx.Done():
		p.log.Error().Err(ctx.Err()).
			Str("topic", topic).
			Msg("sync message sending canceled by context")
		return 0, 0, ctx.Err()

	case <-resultCh:
		if sendErr != nil {
			p.log.Error().Err(sendErr).
				Str("topic", topic).
				Bytes("key", key).
				Bytes("value", value).
				Msg("failed to send sync message")
			return 0, 0, sendErr
		}

		p.log.Debug().
			Str("topic", topic).
			Int32("partition", part).
			Int64("offset", off).
			Msg("sync message sent successfully")

		return part, off, nil
	}
}

func (p *Producer) SendAsync(ctx context.Context, topic string, key, value []byte, headers map[string][]byte) error {
	if p.asyncProducer == nil {
		p.log.Error().
			Str("topic", topic).
			Msg("async producer is not initialized")
		return errors.New("async producer is not initialized")
	}

	msg := p.buildMessage(topic, key, value, headers)

	select {
	case <-ctx.Done():
		p.log.Error().Err(ctx.Err()).
			Str("topic", topic).
			Msg("async message sending canceled by context")
		return ctx.Err()

	case p.asyncProducer.Input() <- msg:
		p.log.Debug().
			Str("topic", topic).
			Bytes("key", key).
			Bytes("value", value).
			Msg("async message sent to producer input channel")
		return nil
	}
}

func (p *Producer) Close() error {
	var err error

	if p.syncProducer != nil {
		if e := p.syncProducer.Close(); e != nil {
			p.log.Error().Err(e).Msg("failed to close sync producer")
			err = e
		} else {
			p.log.Info().Msg("sync producer closed successfully")
		}
	}

	if p.asyncProducer != nil {
		if e := p.asyncProducer.Close(); e != nil {
			p.log.Error().Err(e).Msg("failed to close async producer")
			err = e
		} else {
			p.log.Info().Msg("async producer closed successfully")
		}
		p.closeWg.Wait()
	}

	return err
}

type MessageHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type consumerGroupHandler struct {
	handler             MessageHandler
	consumerWorkerCount int
	log                 *zerolog.Logger
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.log.Info().Msg("consumer group setup completed")
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.log.Info().Msg("consumer group cleanup completed")
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := make(chan *sarama.ConsumerMessage)

	var wg sync.WaitGroup

	for i := range h.consumerWorkerCount {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for msg := range msgCh {
				if err := h.handler(session.Context(), msg); err != nil {
					h.log.Error().Err(err).
						Int("worker id", workerID).
						Str("topic", msg.Topic).
						Int32("partition", msg.Partition).
						Int64("offset", msg.Offset).
						Msg("failed to process message")
					continue
				}
				session.MarkMessage(msg, "")
				h.log.Debug().
					Int("worker id", workerID).
					Str("topic", msg.Topic).
					Int32("partition", msg.Partition).
					Int64("offset", msg.Offset).
					Msg("message processed successfully")
			}
		}(i)
	}

outer:
	for msg := range claim.Messages() {
		select {
		case msgCh <- msg:
		case <-session.Context().Done():
			break outer
		}
	}

	close(msgCh)

	wg.Wait()

	return nil
}

type ConsumerGroup struct {
	group               sarama.ConsumerGroup
	handler             MessageHandler
	topics              []string
	consumerWorkerCount int
	log                 *zerolog.Logger
}

func NewConsumerGroup(k *Kafka, topics []string, handler MessageHandler) (*ConsumerGroup, error) {
	log := k.log.With().Str("component", "consumer group").Logger()

	group, err := sarama.NewConsumerGroup(k.brokers, k.consumerGroup, k.cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create consumer group")
		return nil, err
	}

	log.Info().
		Strs("topics", topics).
		Str("group_id", k.consumerGroup).
		Msg("consumer group created")

	return &ConsumerGroup{
		group:               group,
		handler:             handler,
		topics:              topics,
		consumerWorkerCount: k.consumerWorkerCount,
		log:                 &log,
	}, nil
}

func (c *ConsumerGroup) Start(ctx context.Context) error {
	c.log.Info().
		Strs("topics", c.topics).
		Msg("starting consumer group")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	errCh := make(chan error, 1)

	go func() {
		defer wg.Done()
		for {
			if ctx.Err() != nil {
				return
			}

			handler := &consumerGroupHandler{
				handler:             c.handler,
				consumerWorkerCount: c.consumerWorkerCount,
				log:                 c.log,
			}

			if err := c.group.Consume(ctx, c.topics, handler); err != nil {
				c.log.Error().Err(err).
					Msg("error during consuming messages")
				select {
				case errCh <- err:
				default:
				}
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		wg.Wait()
		c.log.Info().Msg("consumer group context canceled, shutdown complete")
		return nil
	}
}

func (c *ConsumerGroup) Close() error {
	c.log.Info().Msg("closing consumer group")

	if err := c.group.Close(); err != nil {
		c.log.Error().Err(err).Msg("failed to close consumer group")
		return err
	}

	c.log.Info().Msg("consumer group closed successfully")

	return nil
}
