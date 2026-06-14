package review

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(review *model.ProductReview) error
	FindByProductID(productID uint, page, size int) ([]model.ProductReview, int64, error)
	AvgScoreByProductID(productID uint) (float64, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(review *model.ProductReview) error {
	return database.DB.Create(review).Error
}

func (r *repo) FindByProductID(productID uint, page, size int) ([]model.ProductReview, int64, error) {
	var reviews []model.ProductReview
	var total int64
	query := database.DB.Model(&model.ProductReview{}).Where("product_id = ?", productID)
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&reviews).Error
	return reviews, total, err
}

func (r *repo) AvgScoreByProductID(productID uint) (float64, error) {
	var avg float64
	err := database.DB.Model(&model.ProductReview{}).Select("COALESCE(AVG(score), 0)").Where("product_id = ?", productID).Scan(&avg).Error
	return avg, err
}
