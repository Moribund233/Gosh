package order

import (
	"time"

	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByOrderNo(orderNo string) (*model.Order, error)
	ListByUserID(userID uint, status string, page, size int) ([]model.Order, int64, error)
	Update(order *model.Order) error
	UpdateStatus(id uint, status string, version int) error
	CreateItem(item *model.OrderItem) error
	CreateLog(log *model.OrderLog) error
	CreateIdempotency(record *model.IdempotencyRecord) error
	FindIdempotency(key string) (*model.IdempotencyRecord, error)
	FindExpiredOrders(before time.Time) ([]model.Order, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(order *model.Order) error {
	return database.DB.Create(order).Error
}

func (r *repo) FindByID(id uint) (*model.Order, error) {
	var o model.Order
	err := database.DB.Preload("Items").First(&o, id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repo) FindByOrderNo(orderNo string) (*model.Order, error) {
	var o model.Order
	err := database.DB.Where("order_no = ?", orderNo).First(&o).Error
	return &o, err
}

func (r *repo) ListByUserID(userID uint, status string, page, size int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	query := database.DB.Model(&model.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Preload("Items").Order("created_at desc").Offset(offset).Limit(size).Find(&orders).Error
	return orders, total, err
}

func (r *repo) Update(order *model.Order) error {
	return database.DB.Save(order).Error
}

func (r *repo) UpdateStatus(id uint, status string, version int) error {
	return database.DB.Model(&model.Order{}).
		Where("id = ? AND version = ?", id, version).
		Updates(map[string]interface{}{"status": status, "version": version + 1}).Error
}

func (r *repo) CreateItem(item *model.OrderItem) error {
	return database.DB.Create(item).Error
}

func (r *repo) CreateLog(log *model.OrderLog) error {
	return database.DB.Create(log).Error
}

func (r *repo) CreateIdempotency(record *model.IdempotencyRecord) error {
	return database.DB.Create(record).Error
}

func (r *repo) FindIdempotency(key string) (*model.IdempotencyRecord, error) {
	var rec model.IdempotencyRecord
	err := database.DB.Where("key = ?", key).First(&rec).Error
	return &rec, err
}

func (r *repo) FindExpiredOrders(before time.Time) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Where("status = ? AND created_at < ?", model.OrderStatusPendingPayment, before).Find(&orders).Error
	return orders, err
}
