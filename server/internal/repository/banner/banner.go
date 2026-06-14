package banner

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(banner *model.Banner) error
	FindByID(id uint) (*model.Banner, error)
	List(status string) ([]model.Banner, error)
	Update(banner *model.Banner) error
	Delete(id uint) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(banner *model.Banner) error {
	return database.DB.Create(banner).Error
}

func (r *repo) FindByID(id uint) (*model.Banner, error) {
	var b model.Banner
	err := database.DB.First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repo) List(status string) ([]model.Banner, error) {
	var list []model.Banner
	query := database.DB.Model(&model.Banner{}).Order("sort_order asc, id asc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&list).Error
	return list, err
}

func (r *repo) Update(banner *model.Banner) error {
	return database.DB.Save(banner).Error
}

func (r *repo) Delete(id uint) error {
	return database.DB.Delete(&model.Banner{}, id).Error
}
