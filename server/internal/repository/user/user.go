package user

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(user *model.User) error
	FindByPhone(phone string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
	ListByRole(role string, page, size int) ([]model.User, int64, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *repo) FindByPhone(phone string) (*model.User, error) {
	var u model.User
	err := database.DB.Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) FindByID(id uint) (*model.User, error) {
	var u model.User
	err := database.DB.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

func (r *repo) ListByRole(role string, page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := database.DB.Model(&model.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Order("id asc").Offset(offset).Limit(size).Find(&users).Error
	return users, total, err
}
