// Package rabbitmq provides a wrapper for RabbitMQ with separate channels for publishing and consuming, DLX/DLQ, TTL, and worker-based DLQ consumer.
package pkgRabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client represents a RabbitMQ client with separate channels for publishing and consuming.
type Client struct {
	conn            *amqp.Connection
	pubChannel      *amqp.Channel
	consumerChannel *amqp.Channel
	config          Config
}

// NewClient creates a new RabbitMQ client and sets up exchanges and queues.
func NewClient(cfg *Config) (*Client, error) {
	vhost := cfg.VHost
	if vhost == "" {
		vhost = "/"
	}
	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		vhost,
	)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	pubCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open publish channel: %w", err)
	}

	consumerCh, err := conn.Channel()
	if err != nil {
		pubCh.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to open consumer channel: %w", err)
	}

	client := &Client{
		conn:            conn,
		pubChannel:      pubCh,
		consumerChannel: consumerCh,
		config:          *cfg,
	}

	if err := client.setup(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// setup declares exchanges, queues, and bindings
func (c *Client) setup() error {
	// DLX exchange
	if err := c.pubChannel.ExchangeDeclare(
		c.config.DLX,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// DLQ queue
	if _, err := c.pubChannel.QueueDeclare(
		c.config.DLQ,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Bind DLQ to DLX
	if err := c.pubChannel.QueueBind(
		c.config.DLQ,
		"",
		c.config.DLX,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind DLQ to DLX: %w", err)
	}

	// main exchange
	if err := c.pubChannel.ExchangeDeclare(
		c.config.Exchange,
		c.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare main exchange: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": c.config.DLX,
	}

	// main queue with DLX
	if _, err := c.pubChannel.QueueDeclare(
		c.config.Queue,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		return fmt.Errorf("failed to declare main queue: %w", err)
	}

	if err := c.pubChannel.QueueBind(
		c.config.Queue,
		c.config.RoutingKey,
		c.config.Exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind main queue: %w", err)
	}

	return nil
}

// Message represents a message to publish.
type Message struct {
	Body []byte
	TTL  time.Duration
}

// Publish sends a message with dynamic TTL.
func (c *Client) Publish(msg Message) error {
	expiration := fmt.Sprintf("%d", msg.TTL.Milliseconds())
	return c.pubChannel.Publish(
		c.config.Exchange,
		c.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			Expiration:  expiration,
		},
	)
}

// PublishStructWithTTL automatically serializes the structure to JSON and publishes the message.
func (c *Client) PublishStructWithTTL(data interface{}, ttl time.Duration) error {
	if err := c.ensurePubChannel(); err != nil {
		return fmt.Errorf("failed to ensure publish channel: %w", err)
	}

	if ttl <= 0 {
		return fmt.Errorf("invalid TTL: must be positive, got %v", ttl)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	queueName := fmt.Sprintf("msg_%d", time.Now().UnixNano())

	tempQueue, err := c.DeclareTempQueue(queueName, ttl)
	if err != nil {
		return fmt.Errorf("failed to declare temp queue: %w", err)
	}

	if err := c.pubChannel.QueueBind(
		tempQueue,
		tempQueue,
		c.config.Exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind temp queue %q: %w", tempQueue, err)
	}

	expiration := fmt.Sprintf("%d", ttl.Milliseconds())
	err = c.pubChannel.Publish(
		c.config.Exchange,
		tempQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  expiration,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to temp queue %q: %w", tempQueue, err)
	}

	return nil
}

func (c *Client) ensurePubChannel() error {
	if c.pubChannel != nil && !c.pubChannel.IsClosed() {
		return nil
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to reopen publish channel: %w", err)
	}

	c.pubChannel = ch
	return nil
}

// DeclareTempQueue creates a temporary queue with a specific TTL and auto-delete flag.
func (c *Client) DeclareTempQueue(name string, ttl time.Duration) (string, error) {
	ms := ttl.Milliseconds()
	if ms <= 0 {
		return "", fmt.Errorf("invalid TTL: must be positive, got %d ms", ms)
	}

	exp := ms + 5000
	if exp <= 0 {
		return "", fmt.Errorf("invalid x-expires: overflow or negative value (%d)", exp)
	}

	args := amqp.Table{
		"x-message-ttl":          int64(ms),
		"x-expires":              int64(exp),
		"x-dead-letter-exchange": c.config.DLX,
	}

	q, err := c.pubChannel.QueueDeclare(
		name,
		false,
		true,
		false,
		false,
		args,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare temp queue %q: %w", name, err)
	}

	return q.Name, nil
}

// Ack acknowledges a message.
func (c *Client) Ack(msg amqp.Delivery) error {
	return msg.Ack(false)
}

// Nack negatively acknowledges a message.
func (c *Client) Nack(msg amqp.Delivery) error {
	return msg.Nack(false, false)
}

// Close closes both channels and the connection.
func (c *Client) Close() error {
	if c.pubChannel != nil {
		c.pubChannel.Close()
	}
	if c.consumerChannel != nil {
		c.consumerChannel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ConsumeDLQWithWorkers consumes messages from DLQ with a worker pool and context for graceful shutdown.
func (c *Client) ConsumeDLQWithWorkers(ctx context.Context, workerCount int, handler func(msg amqp.Delivery)) error {
	if workerCount <= 0 {
		workerCount = 1
	}

	const consumerTag = "dlq-consumer"

	msgs, err := c.consumerChannel.Consume(
		c.config.DLQ,
		consumerTag,
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume from DLQ: %w", err)
	}

	queue := make(chan amqp.Delivery, workerCount*2)

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case msg, ok := <-queue:
					if !ok {
						return
					}
					handler(msg)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		for msg := range msgs {
			select {
			case queue <- msg:
			case <-ctx.Done():
				close(queue)
				return
			}
		}
		close(queue)
	}()

	go func() {
		<-ctx.Done()
		_ = c.CancelConsumer(consumerTag) // gracefully cancel consumer
		_ = c.consumerChannel.Close()
	}()

	return nil
}

// CancelConsumer gracefully cancels a consumer by tag.
func (c *Client) CancelConsumer(consumerTag string) error {
	if c.consumerChannel == nil {
		return fmt.Errorf("consumer channel is nil")
	}
	if err := c.consumerChannel.Cancel(consumerTag, false); err != nil {
		return fmt.Errorf("failed to cancel consumer: %w", err)
	}
	return nil
}
