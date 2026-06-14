package cart

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(cart *model.Cart) error
	FindByID(id, userID uint) (*model.Cart, error)
	FindByUserID(userID uint) ([]model.Cart, error)
	FindByUserAndSKU(userID, skuID uint) (*model.Cart, error)
	FindSelectedByUserID(userID uint) ([]model.Cart, error)
	Update(cart *model.Cart) error
	Delete(id, userID uint) error
	DeleteBySKUIDs(userID uint, skuIDs []uint) error
	DeleteByUserID(userID uint) error
	DeleteSelectedByUserID(userID uint) error
	Count(userID uint) (int64, error)
	UpdateSelected(userID, skuID uint, selected bool) error
	UpdateAllSelected(userID uint, selected bool) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(cart *model.Cart) error {
	return database.DB.Create(cart).Error
}

func (r *repo) FindByID(id, userID uint) (*model.Cart, error) {
	var c model.Cart
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repo) FindByUserID(userID uint) ([]model.Cart, error) {
	var carts []model.Cart
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&carts).Error
	return carts, err
}

func (r *repo) FindByUserAndSKU(userID, skuID uint) (*model.Cart, error) {
	var c model.Cart
	err := database.DB.Where("user_id = ? AND sku_id = ?", userID, skuID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repo) FindSelectedByUserID(userID uint) ([]model.Cart, error) {
	var carts []model.Cart
	err := database.DB.Where("user_id = ? AND selected = ?", userID, true).Find(&carts).Error
	return carts, err
}

func (r *repo) DeleteSelectedByUserID(userID uint) error {
	return database.DB.Where("user_id = ? AND selected = ?", userID, true).Delete(&model.Cart{}).Error
}

func (r *repo) Update(cart *model.Cart) error {
	return database.DB.Save(cart).Error
}

func (r *repo) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Cart{}).Error
}

func (r *repo) DeleteBySKUIDs(userID uint, skuIDs []uint) error {
	if len(skuIDs) == 0 {
		return nil
	}
	return database.DB.Where("user_id = ? AND sku_id IN ?", userID, skuIDs).Delete(&model.Cart{}).Error
}

func (r *repo) DeleteByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

func (r *repo) Count(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Cart{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *repo) UpdateSelected(userID, skuID uint, selected bool) error {
	return database.DB.Model(&model.Cart{}).
		Where("user_id = ? AND sku_id = ?", userID, skuID).
		Update("selected", selected).Error
}

func (r *repo) UpdateAllSelected(userID uint, selected bool) error {
	return database.DB.Model(&model.Cart{}).
		Where("user_id = ?", userID).
		Update("selected", selected).Error
}
