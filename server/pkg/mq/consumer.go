package mq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Handler func(ctx context.Context, msg []byte) error

type Consumer struct {
	queue    string
	handler  Handler
	logger   *zap.Logger
	stopChan chan struct{}
}

func NewConsumer(queue string, handler Handler, logger *zap.Logger) *Consumer {
	return &Consumer{
		queue:    queue,
		handler:  handler,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (c *Consumer) Start() {
	if DefaultConn == nil || DefaultConn.channel == nil {
		c.logger.Warn("mq not available, consumer not started", zap.String("queue", c.queue))
		return
	}

	msgs, err := DefaultConn.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("consume failed", zap.String("queue", c.queue), zap.Error(err))
		return
	}

	c.logger.Info("consumer started", zap.String("queue", c.queue))

	go func() {
		for {
			select {
			case <-c.stopChan:
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				c.process(d)
			}
		}
	}()
}

func (c *Consumer) Stop() {
	close(c.stopChan)
}

func (c *Consumer) process(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.handler(ctx, d.Body); err != nil {
		c.logger.Error("message processing failed",
			zap.String("queue", c.queue),
			zap.String("body", string(d.Body)),
			zap.Error(err),
		)
		d.Nack(false, true)
		return
	}

	d.Ack(false)
}

func ParsePayload[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
