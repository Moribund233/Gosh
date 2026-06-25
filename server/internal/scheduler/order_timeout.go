package scheduler

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
)

type Scheduler struct {
	timeout time.Duration
	stopCh  chan struct{}
}

func New(timeoutMinutes int) *Scheduler {
	return &Scheduler{
		timeout: time.Duration(timeoutMinutes) * time.Minute,
		stopCh:  make(chan struct{}),
	}
}

func (s *Scheduler) Start(log *zap.Logger) {
	log.Info("order timeout scheduler started", zap.Duration("interval", 1*time.Minute), zap.Duration("timeout", s.timeout))
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cancelExpiredOrders(log)
			case <-s.stopCh:
				log.Info("order timeout scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) cancelExpiredOrders(log *zap.Logger) {
	before := time.Now().Add(-s.timeout)
	var orders []model.Order

	if err := database.DB.Where("status = ? AND created_at < ?", model.OrderStatusPendingPayment, before).Find(&orders).Error; err != nil {
		log.Error("query expired orders failed", zap.Error(err))
		return
	}

	if len(orders) == 0 {
		return
	}

	cancelled := 0
	for _, order := range orders {
		txCtl, cancel := database.WithTimeout(database.DefaultTimeout)
		err := txCtl.Transaction(func(tx *gorm.DB) error {
			var items []model.OrderItem
			if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
				return err
			}

			res := tx.Model(&model.Order{}).
				Where("id = ? AND version = ?", order.ID, order.Version).
				Updates(map[string]interface{}{
					"status":       model.OrderStatusCancelled,
					"cancelled_at": time.Now(),
					"version":      order.Version + 1,
				})
			if res.RowsAffected == 0 {
				return nil
			}

			for _, item := range items {
				if err := tx.Model(&model.ProductSKU{}).
					Where("id = ?", item.SKUID).
					Updates(map[string]interface{}{
						"stock":   gorm.Expr("stock + ?", item.Quantity),
						"version": gorm.Expr("version + 1"),
					}).Error; err != nil {
					return err
				}
			}

			return tx.Create(&model.OrderLog{
				OrderID:    order.ID,
				FromStatus: model.OrderStatusPendingPayment,
				ToStatus:   model.OrderStatusCancelled,
				Operator:   "system",
				Note:       "超时未支付，系统自动取消",
			}).Error
		})
		cancel()

		if err != nil {
			log.Error("cancel expired order failed", zap.Uint("order_id", order.ID), zap.Error(err))
			continue
		}
		cancelled++
	}

	if cancelled > 0 {
		log.Info("auto-cancelled expired orders", zap.Int("count", cancelled))
	}
}
