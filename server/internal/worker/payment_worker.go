package worker

import (
	"context"

	"go.uber.org/zap"
	"gosh/pkg/mq"
	paymentSvc "gosh/internal/service/payment"
)

type PaymentCallbackEvent struct {
	Method string `json:"method"`
	Body   []byte `json:"body"`
}

func NewPaymentWorker(logger *zap.Logger) *mq.Consumer {
	return mq.NewConsumer(mq.QueuePaymentCallback, func(ctx context.Context, body []byte) error {
		event, err := mq.ParsePayload[PaymentCallbackEvent](body)
		if err != nil {
			return err
		}

		svc := paymentSvc.New()
		if err := svc.ProcessCallback(event.Method, event.Body); err != nil {
			logger.Error("payment callback processing failed",
				zap.String("method", event.Method),
				zap.Error(err),
			)
			return err
		}

		logger.Info("payment callback processed",
			zap.String("method", event.Method),
		)
		return nil
	}, logger)
}
