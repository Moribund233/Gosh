package browse_history

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(history *model.BrowseHistory) error
	FindByUserID(userID uint, page, size int) ([]model.BrowseHistory, int64, error)
	CountByUserID(userID uint) (int64, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(history *model.BrowseHistory) error {
	return database.DB.Create(history).Error
}

func (r *repo) FindByUserID(userID uint, page, size int) ([]model.BrowseHistory, int64, error) {
	var list []model.BrowseHistory
	var total int64
	query := database.DB.Model(&model.BrowseHistory{}).Where("user_id = ?", userID)
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *repo) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.BrowseHistory{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
