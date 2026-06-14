package product

import (
	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(product *model.Product, skus []model.ProductSKU) error
	FindByID(id uint) (*model.Product, error)
	FindSKUByID(skuID uint) (*model.ProductSKU, error)
	FindSKUsByProductID(productID uint) ([]model.ProductSKU, error)
	List(categoryID uint, tag, keyword, sort, status string, page, size int) ([]model.Product, int64, error)
	Update(product *model.Product) error
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(product *model.Product, skus []model.ProductSKU) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		for i := range skus {
			skus[i].ProductID = product.ID
		}
		if len(skus) > 0 {
			if err := tx.Create(&skus).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repo) FindByID(id uint) (*model.Product, error) {
	var p model.Product
	err := database.DB.Preload("SKUs").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) FindSKUByID(skuID uint) (*model.ProductSKU, error) {
	var sku model.ProductSKU
	err := database.DB.First(&sku, skuID).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *repo) FindSKUsByProductID(productID uint) ([]model.ProductSKU, error) {
	var skus []model.ProductSKU
	err := database.DB.Where("product_id = ?", productID).Find(&skus).Error
	return skus, err
}

func (r *repo) List(categoryID uint, tag, keyword, sort, status string, page, size int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	query := database.DB.Model(&model.Product{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR subtitle LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", model.ProductStatusOn)
	}
	query.Count(&total)
	offset := (page - 1) * size
	switch sort {
	case "price_asc":
		query = query.Order("price asc")
	case "price_desc":
		query = query.Order("price desc")
	case "sales":
		query = query.Order("sales desc")
	case "newest":
		query = query.Order("created_at desc")
	default:
		query = query.Order("created_at desc")
	}
	err := query.Offset(offset).Limit(size).Find(&products).Error
	return products, total, err
}

func (r *repo) Update(product *model.Product) error {
	return database.DB.Save(product).Error
}

func (r *repo) UpdateStatus(id uint, status string) error {
	return database.DB.Model(&model.Product{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repo) Delete(id uint) error {
	return database.DB.Delete(&model.Product{}, id).Error
}
