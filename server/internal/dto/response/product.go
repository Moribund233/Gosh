package response

import (
	"gosh/internal/model"
)

type ProductResponse struct {
	ID            uint              `json:"id"`
	CategoryID    uint              `json:"category_id"`
	Name          string            `json:"name"`
	Subtitle      string            `json:"subtitle"`
	Brand         string            `json:"brand"`
	Price         int64             `json:"price"`
	OriginalPrice int64             `json:"original_price"`
	Sales         int64             `json:"sales"`
	Tags          string            `json:"tags"`
	Images        []string          `json:"images"`
	Status        string            `json:"status"`
	Description   string            `json:"description,omitempty"`
	Origin        string            `json:"origin"`
	ShelfLife     string            `json:"shelf_life"`
	IsNew         bool              `json:"is_new"`
	IsHot         bool              `json:"is_hot"`
	IsFeatured    bool              `json:"is_featured"`
	CreatedAt     string            `json:"created_at"`
	SKUs          []SKUResponse     `json:"skus,omitempty"`
}

type SKUResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Stock     int    `json:"stock"`
}

type ReviewResponse struct {
	ID        uint     `json:"id"`
	ProductID uint     `json:"product_id"`
	UserID    uint     `json:"user_id"`
	Score     int      `json:"score"`
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

func ToProductResponse(p *model.Product) ProductResponse {
	return ProductResponse{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		Name:          p.Name,
		Subtitle:      p.Subtitle,
		Brand:         p.Brand,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Sales:         p.Sales,
		Tags:          p.Tags,
		Images:        p.Images,
		Status:        p.Status,
		Description:   p.Description,
		Origin:        p.Origin,
		ShelfLife:     p.ShelfLife,
		IsNew:         p.IsNew,
		IsHot:         p.IsHot,
		IsFeatured:    p.IsFeatured,
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToProductList(products []model.Product) []ProductResponse {
	list := make([]ProductResponse, len(products))
	for i, p := range products {
		list[i] = ToProductResponse(&p)
	}
	return list
}

func ToSKUResponse(sku *model.ProductSKU) SKUResponse {
	return SKUResponse{
		ID:        sku.ID,
		ProductID: sku.ProductID,
		Name:      sku.Name,
		Price:     sku.Price,
		Stock:     sku.Stock,
	}
}

func ToSKUList(skus []model.ProductSKU) []SKUResponse {
	list := make([]SKUResponse, len(skus))
	for i, s := range skus {
		list[i] = ToSKUResponse(&s)
	}
	return list
}

func ToReviewResponse(r *model.ProductReview) ReviewResponse {
	return ReviewResponse{
		ID:        r.ID,
		ProductID: r.ProductID,
		UserID:    r.UserID,
		Score:     r.Score,
		Content:   r.Content,
		Images:    r.Images,
		Nickname:  r.Nickname,
		Avatar:    r.Avatar,
		CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToReviewList(reviews []model.ProductReview) []ReviewResponse {
	list := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		list[i] = ToReviewResponse(&r)
	}
	return list
}
