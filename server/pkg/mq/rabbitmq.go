package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gosh/internal/config"
)

const (
	ExchangeGosh     = "gosh"
	ExchangeGoshKind = "topic"

	QueuePointAward    = "gosh.point.award"
	QueuePaymentCallback = "gosh.payment.callback"

	RoutingKeyOrderPaid     = "order.paid"
	RoutingKeyPaymentCallback = "payment.callback"
)

var DefaultConn *Connection

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
}

func Init(cfg config.RabbitMQConfig, logger *zap.Logger) error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost,
	)
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("rabbitmq dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("rabbitmq channel failed: %w", err)
	}

	c := &Connection{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}

	if err := c.declare(); err != nil {
		conn.Close()
		return fmt.Errorf("rabbitmq declare failed: %w", err)
	}

	DefaultConn = c
	logger.Info("rabbitmq connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)
	return nil
}

func (c *Connection) declare() error {
	if err := c.channel.ExchangeDeclare(
		ExchangeGosh,
		ExchangeGoshKind,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	queues := []string{QueuePointAward, QueuePaymentCallback}
	for _, q := range queues {
		if _, err := c.channel.QueueDeclare(
			q,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("queue declare %s failed: %w", q, err)
		}
	}

	bindings := map[string]string{
		QueuePointAward:    RoutingKeyOrderPaid,
		QueuePaymentCallback: RoutingKeyPaymentCallback,
	}
	for queue, key := range bindings {
		if err := c.channel.QueueBind(queue, key, ExchangeGosh, false, nil); err != nil {
			return fmt.Errorf("queue bind %s -> %s failed: %w", queue, key, err)
		}
	}

	return nil
}

func (c *Connection) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}
