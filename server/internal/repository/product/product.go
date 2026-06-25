package product

import (
	"strings"

	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(product *model.Product, skus []model.ProductSKU) error
	FindByID(id uint) (*model.Product, error)
	FindByIDs(ids []uint) ([]model.Product, error)
	FindSKUByID(skuID uint) (*model.ProductSKU, error)
	FindSKUsByIDs(skuIDs []uint) ([]model.ProductSKU, error)
	FindSKUsByProductID(productID uint) ([]model.ProductSKU, error)
	IsProductOnline(productID uint) (bool, error)
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

func (r *repo) FindSKUsByIDs(skuIDs []uint) ([]model.ProductSKU, error) {
	if len(skuIDs) == 0 {
		return nil, nil
	}
	var skus []model.ProductSKU
	err := database.DB.Where("id IN ?", skuIDs).Find(&skus).Error
	return skus, err
}

func (r *repo) FindByIDs(ids []uint) ([]model.Product, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var products []model.Product
	err := database.DB.Select("id, name, images, status").Where("id IN ?", ids).Find(&products).Error
	return products, err
}

func (r *repo) IsProductOnline(productID uint) (bool, error) {
	var status string
	err := database.DB.Model(&model.Product{}).Select("status").Where("id = ?", productID).Scan(&status).Error
	if err != nil {
		return false, err
	}
	return status == model.ProductStatusOn, nil
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
		query = keywordFilter(query, keyword)
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

func keywordFilter(query *gorm.DB, keyword string) *gorm.DB {
	driver := database.DB.Dialector.Name()
	if driver == "postgres" {
		cleaned := sanitizeFTS(keyword)
		if cleaned == "" {
			return query.Where("1 = 0")
		}
		return query.Where(
			"to_tsvector('simple', coalesce(name,'') || ' ' || coalesce(subtitle,'')) @@ plainto_tsquery('simple', ?)",
			cleaned,
		)
	}
	return query.Where("name LIKE ? OR subtitle LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
}

func sanitizeFTS(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == ' ' || r == '-' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r >= 0x80 {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}
