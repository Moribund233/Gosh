package point

import (
	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
	repo "gosh/internal/repository/point"
)

type Service interface {
	GetBalance(userID uint) (int, error)
	ListLogs(userID uint, page, size int) ([]model.PointLog, int64, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) GetBalance(userID uint) (int, error) {
	return s.repo.GetUserPoints(userID)
}

func (s *service) ListLogs(userID uint, page, size int) ([]model.PointLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	return s.repo.ListByUserID(userID, page, size)
}

func EarnPoints(userID uint, orderID uint, amount int64) error {
	points := int(amount / 100)
	if points <= 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			Update("points", gorm.Expr("points + ?", points)).Error; err != nil {
			return err
		}

		var user model.User
		tx.Select("points").First(&user, userID)

		return tx.Create(&model.PointLog{
			UserID:  userID,
			Type:    model.PointTypeEarn,
			Amount:  points,
			Balance: user.Points,
			OrderID: &orderID,
			Note:    "下单赠送积分",
		}).Error
	})
}
