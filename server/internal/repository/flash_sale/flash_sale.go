package flash_sale

import (
	"time"

	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	FindActive() ([]model.FlashSale, error)
	FindByProductID(productID uint) (*model.FlashSale, error)
	DecrementStock(id uint, version int) bool
	Create(fs *model.FlashSale) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) FindActive() ([]model.FlashSale, error) {
	var list []model.FlashSale
	now := time.Now()
	err := database.DB.Where("status = ? AND start_at <= ? AND end_at > ?",
		model.FlashSaleStatusActive, now, now).
		Find(&list).Error
	return list, err
}

func (r *repo) FindByProductID(productID uint) (*model.FlashSale, error) {
	var fs model.FlashSale
	now := time.Now()
	err := database.DB.Where("product_id = ? AND status = ? AND start_at <= ? AND end_at > ?",
		productID, model.FlashSaleStatusActive, now, now).
		First(&fs).Error
	return &fs, err
}

func (r *repo) DecrementStock(id uint, version int) bool {
	res := database.DB.Model(&model.FlashSale{}).
		Where("id = ? AND flash_stock > 0 AND version = ?", id, version).
		Updates(map[string]interface{}{
			"flash_stock": gorm.Expr("flash_stock - 1"),
			"version":     version + 1,
		})
	return res.RowsAffected > 0
}

func (r *repo) Create(fs *model.FlashSale) error {
	return database.DB.Create(fs).Error
}
