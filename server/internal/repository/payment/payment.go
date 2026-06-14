package payment

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(payment *model.Payment) error
	FindByID(id uint) (*model.Payment, error)
	FindByOrderNo(orderNo string) (*model.Payment, error)
	FindByTransactionNo(txNo string) (*model.Payment, error)
	Update(payment *model.Payment) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(payment *model.Payment) error {
	return database.DB.Create(payment).Error
}

func (r *repo) FindByID(id uint) (*model.Payment, error) {
	var p model.Payment
	err := database.DB.First(&p, id).Error
	return &p, err
}

func (r *repo) FindByOrderNo(orderNo string) (*model.Payment, error) {
	var p model.Payment
	err := database.DB.Where("order_no = ?", orderNo).First(&p).Error
	return &p, err
}

func (r *repo) FindByTransactionNo(txNo string) (*model.Payment, error) {
	var p model.Payment
	err := database.DB.Where("transaction_no = ?", txNo).First(&p).Error
	return &p, err
}

func (r *repo) Update(payment *model.Payment) error {
	return database.DB.Save(payment).Error
}
