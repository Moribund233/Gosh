package category

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(cat *model.Category) error
	FindByID(id uint) (*model.Category, error)
	FindAll() ([]model.Category, error)
	FindByParentID(parentID uint) ([]model.Category, error)
	Update(cat *model.Category) error
	Delete(id uint) error
	CountByParentID(parentID uint) (int64, error)
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(cat *model.Category) error {
	return database.DB.Create(cat).Error
}

func (r *repo) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	err := database.DB.First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repo) FindAll() ([]model.Category, error) {
	var cats []model.Category
	err := database.DB.Order("sort_order asc, id asc").Find(&cats).Error
	return cats, err
}

func (r *repo) FindByParentID(parentID uint) ([]model.Category, error) {
	var cats []model.Category
	err := database.DB.Where("parent_id = ?", parentID).Order("sort_order asc, id asc").Find(&cats).Error
	return cats, err
}

func (r *repo) Update(cat *model.Category) error {
	return database.DB.Save(cat).Error
}

func (r *repo) Delete(id uint) error {
	return database.DB.Delete(&model.Category{}, id).Error
}

func (r *repo) CountByParentID(parentID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Category{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}
