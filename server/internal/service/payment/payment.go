package payment

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gosh/internal/database"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	orderRepo "gosh/internal/repository/order"
	paymentRepo "gosh/internal/repository/payment"
	paymentPkg "gosh/internal/pkg/payment"
	pointSvc "gosh/internal/service/point"
	"gosh/pkg/mq"
)

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderNotBelongToUser = errors.New("order does not belong to user")
	ErrInvalidOrderStatus   = errors.New("invalid order status for payment")
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrRefundNotImplemented = errors.New("refund not yet implemented")
	ErrRefundExists         = errors.New("refund already processed for this payment")
)

type Service interface {
	GetMethods() []response.PaymentMethodResponse
	Pay(userID uint, req *request.PayRequest) (*response.PaymentResponse, error)
	ProcessCallback(method string, body []byte) error
	GetStatus(orderNo string) (*response.PaymentResponse, error)
	Refund(userID uint, req *request.RefundRequest) error
}

type service struct {
	paymentRepo paymentRepo.Repository
	orderRepo   orderRepo.Repository
	provider    paymentPkg.Provider
}

func New() Service {
	return &service{
		paymentRepo: paymentRepo.New(),
		orderRepo:   orderRepo.New(),
		provider:    paymentPkg.NewMockProvider(),
	}
}

func (s *service) GetMethods() []response.PaymentMethodResponse {
	return []response.PaymentMethodResponse{
		{Method: model.PaymentMethodMock, Name: "模拟支付"},
		{Method: model.PaymentMethodWechat, Name: "微信支付"},
		{Method: model.PaymentMethodAlipay, Name: "支付宝"},
	}
}

func (s *service) Pay(userID uint, req *request.PayRequest) (*response.PaymentResponse, error) {
	order, err := s.orderRepo.FindByOrderNo(req.OrderNo)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingPayment {
		return nil, ErrInvalidOrderStatus
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		payment := &model.Payment{
			OrderID:   order.ID,
			OrderNo:   order.OrderNo,
			Method:    req.Method,
			PayAmount: order.PayAmount,
			Status:    model.PaymentStatusPending,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		result, err := s.provider.CreatePayment(order, req.Method)
		if err != nil {
			return err
		}

		callbackResult, err := s.provider.ProcessCallback(req.Method, []byte(result.RawData))
		if err != nil {
			return err
		}
		if !callbackResult.SignOk {
			return fmt.Errorf("payment signature verification failed")
		}

		now := time.Now()
		payment.TransactionNo = callbackResult.TransactionNo
		payment.NotifyRaw = callbackResult.RawData
		payment.NotifySignOk = callbackResult.SignOk
		payment.Status = model.PaymentStatusSuccess
		payment.PaidAt = &now
		if err := tx.Save(payment).Error; err != nil {
			return err
		}

		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":   model.OrderStatusPendingDelivery,
				"paid_at":  now,
				"version":  order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}

		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingPayment,
			ToStatus:   model.OrderStatusPendingDelivery,
			Operator:   fmt.Sprintf("user:%d", userID),
			Note:       fmt.Sprintf("支付成功 (%s)", req.Method),
		}).Error
	})

	if err != nil {
		return nil, err
	}

	if mq.DefaultConn != nil {
		mq.PublishEvent(mq.RoutingKeyOrderPaid, map[string]interface{}{
			"user_id":    userID,
			"order_id":   order.ID,
			"pay_amount": order.PayAmount,
		})
	} else {
		if err := pointSvc.EarnPoints(userID, order.ID, order.PayAmount); err != nil {
			return nil, err
		}
	}

	fullPayment, err := s.paymentRepo.FindByOrderNo(req.OrderNo)
	if err != nil {
		return nil, err
	}
	resp := response.ToPaymentResponse(fullPayment)
	return &resp, nil
}

func (s *service) ProcessCallback(method string, body []byte) error {
	result, err := s.provider.ProcessCallback(method, body)
	if err != nil {
		return err
	}
	if !result.SignOk {
		return fmt.Errorf("signature verification failed")
	}

	existing, err := s.paymentRepo.FindByTransactionNo(result.TransactionNo)
	if err == nil && existing != nil {
		return nil
	}

	order, err := s.orderRepo.FindByOrderNo(result.OrderNo)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.PayAmount != result.Amount {
		return fmt.Errorf("amount mismatch: expected %d, got %d", order.PayAmount, result.Amount)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return nil
	}

	now := time.Now()
	return database.DB.Transaction(func(tx *gorm.DB) error {
		payment := &model.Payment{
			OrderID:       order.ID,
			OrderNo:       order.OrderNo,
			Method:        method,
			PayAmount:     result.Amount,
			Status:        model.PaymentStatusSuccess,
			TransactionNo: result.TransactionNo,
			PaidAt:        &now,
			NotifyRaw:     result.RawData,
			NotifySignOk:  result.SignOk,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":  model.OrderStatusPendingDelivery,
				"paid_at": now,
				"version": order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}

		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: model.OrderStatusPendingPayment,
			ToStatus:   model.OrderStatusPendingDelivery,
			Operator:   "system",
			Note:       fmt.Sprintf("支付回调成功 (%s)", method),
		}).Error
	})
}

func (s *service) GetStatus(orderNo string) (*response.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByOrderNo(orderNo)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	resp := response.ToPaymentResponse(payment)
	return &resp, nil
}

func (s *service) Refund(userID uint, req *request.RefundRequest) error {
	payment, err := s.paymentRepo.FindByTransactionNo(req.TransactionNo)
	if err != nil {
		return ErrPaymentNotFound
	}
	if payment.Status == model.PaymentStatusRefunded {
		return ErrRefundExists
	}
	if payment.Status != model.PaymentStatusSuccess {
		return fmt.Errorf("payment not in success status")
	}

	order, err := s.orderRepo.FindByOrderNo(payment.OrderNo)
	if err != nil {
		return ErrOrderNotFound
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Payment{}).
			Where("id = ?", payment.ID).
			Update("status", model.PaymentStatusRefunded).Error; err != nil {
			return err
		}

		var items []model.OrderItem
		if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Model(&model.ProductSKU{}).
				Where("id = ?", item.SKUID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":        model.OrderStatusCancelled,
				"cancelled_at":  now,
				"cancel_reason": fmt.Sprintf("退款: %s", req.Reason),
				"version":       order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}

		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: order.Status,
			ToStatus:   model.OrderStatusCancelled,
			Operator:   fmt.Sprintf("admin:%d", userID),
			Note:       fmt.Sprintf("退款处理: %s (交易号: %s)", req.Reason, req.TransactionNo),
		}).Error
	})
}
