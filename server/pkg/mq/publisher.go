package mq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func Publish(ctx context.Context, routingKey string, payload interface{}) error {
	if DefaultConn == nil {
		return nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return DefaultConn.channel.PublishWithContext(ctx,
		ExchangeGosh,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func PublishEvent(routingKey string, payload interface{}) {
	if DefaultConn == nil {
		return
	}
	if err := Publish(context.Background(), routingKey, payload); err != nil && DefaultConn.logger != nil {
		DefaultConn.logger.Warn("mq publish failed",
			zap.String("routing_key", routingKey),
			zap.Error(err),
		)
	}
}
