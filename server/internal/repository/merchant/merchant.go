package merchant

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(app *model.MerchantApplication) error
	FindByID(id uint) (*model.MerchantApplication, error)
	FindByUserID(userID uint) (*model.MerchantApplication, error)
	List(status string, page, size int) ([]model.MerchantApplication, int64, error)
	Update(app *model.MerchantApplication) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(app *model.MerchantApplication) error {
	return database.DB.Create(app).Error
}

func (r *repo) FindByID(id uint) (*model.MerchantApplication, error) {
	var a model.MerchantApplication
	err := database.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) FindByUserID(userID uint) (*model.MerchantApplication, error) {
	var a model.MerchantApplication
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) List(status string, page, size int) ([]model.MerchantApplication, int64, error) {
	var apps []model.MerchantApplication
	var total int64
	query := database.DB.Model(&model.MerchantApplication{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&apps).Error
	return apps, total, err
}

func (r *repo) Update(app *model.MerchantApplication) error {
	return database.DB.Save(app).Error
}
