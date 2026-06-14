package address

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(addr *model.Address) error
	FindByID(id uint) (*model.Address, error)
	FindByUserID(userID uint) ([]model.Address, error)
	FindDefaultByUserID(userID uint) (*model.Address, error)
	Update(addr *model.Address) error
	Delete(id uint, userID uint) error
	ResetDefault(userID uint) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(addr *model.Address) error {
	return database.DB.Create(addr).Error
}

func (r *repo) FindByID(id uint) (*model.Address, error) {
	var a model.Address
	err := database.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) FindByUserID(userID uint) ([]model.Address, error) {
	var addrs []model.Address
	err := database.DB.Where("user_id = ?", userID).Order("is_default desc, created_at desc").Find(&addrs).Error
	return addrs, err
}

func (r *repo) FindDefaultByUserID(userID uint) (*model.Address, error) {
	var a model.Address
	err := database.DB.Where("user_id = ? AND is_default = ?", userID, true).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) Update(addr *model.Address) error {
	return database.DB.Save(addr).Error
}

func (r *repo) Delete(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Address{}).Error
}

func (r *repo) ResetDefault(userID uint) error {
	return database.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}
