package worker

import (
	"context"

	"go.uber.org/zap"
	"gosh/internal/service/point"
	"gosh/pkg/mq"
)

type OrderPaidEvent struct {
	UserID    uint `json:"user_id"`
	OrderID   uint `json:"order_id"`
	PayAmount int64 `json:"pay_amount"`
}

func NewPointWorker(logger *zap.Logger) *mq.Consumer {
	return mq.NewConsumer(mq.QueuePointAward, func(ctx context.Context, body []byte) error {
		event, err := mq.ParsePayload[OrderPaidEvent](body)
		if err != nil {
			return err
		}

		if err := point.EarnPoints(event.UserID, event.OrderID, event.PayAmount); err != nil {
			logger.Error("earn points failed",
				zap.Uint("user_id", event.UserID),
				zap.Uint("order_id", event.OrderID),
				zap.Error(err),
			)
			return err
		}

		logger.Info("points awarded",
			zap.Uint("user_id", event.UserID),
			zap.Uint("order_id", event.OrderID),
			zap.Int64("amount", event.PayAmount),
		)
		return nil
	}, logger)
}
