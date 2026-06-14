package favorite

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(fav *model.Favorite) error
	Delete(userID, productID uint) error
	FindByUserID(userID uint, page, size int) ([]model.Favorite, int64, error)
	Exists(userID, productID uint) (bool, error)
	CountByUserID(userID uint) (int64, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(fav *model.Favorite) error {
	return database.DB.Create(fav).Error
}

func (r *repo) Delete(userID, productID uint) error {
	return database.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Favorite{}).Error
}

func (r *repo) FindByUserID(userID uint, page, size int) ([]model.Favorite, int64, error) {
	var favs []model.Favorite
	var total int64
	query := database.DB.Model(&model.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&favs).Error
	return favs, total, err
}

func (r *repo) Exists(userID, productID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}

func (r *repo) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
