// Package rabbitmq provides a wrapper for RabbitMQ with separate channels for publishing and consuming, DLX/DLQ, TTL, and worker-based DLQ consumer.
package rabbitmq

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
func NewClient(cfg Config) (*Client, error) {
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
		config:          cfg,
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

// PublishStruct automatically serializes the structure to JSON and publishes the message.
func (c *Client) PublishStructWithTTL(data interface{}, ttl time.Duration) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}
	return c.Publish(Message{Body: body, TTL: ttl})
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

	msgs, err := c.consumerChannel.Consume(
		c.config.DLQ,
		"",
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

	return nil
}
