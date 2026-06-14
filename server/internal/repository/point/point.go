package point

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(log *model.PointLog) error
	ListByUserID(userID uint, page, size int) ([]model.PointLog, int64, error)
	GetUserPoints(userID uint) (int, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(log *model.PointLog) error {
	return database.DB.Create(log).Error
}

func (r *repo) ListByUserID(userID uint, page, size int) ([]model.PointLog, int64, error) {
	var logs []model.PointLog
	var total int64
	query := database.DB.Model(&model.PointLog{}).Where("user_id = ?", userID)
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Order("created_at desc").Offset(offset).Limit(size).Find(&logs).Error
	return logs, total, err
}

func (r *repo) GetUserPoints(userID uint) (int, error) {
	var user model.User
	err := database.DB.Select("points").First(&user, userID).Error
	if err != nil {
		return 0, err
	}
	return user.Points, nil
}
