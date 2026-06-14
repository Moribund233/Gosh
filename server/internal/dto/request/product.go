package request

type CreateProductRequest struct {
	CategoryID    uint           `json:"category_id" binding:"required"`
	Name          string         `json:"name" binding:"required,max=128"`
	Subtitle      string         `json:"subtitle" binding:"omitempty,max=256"`
	Brand         string         `json:"brand" binding:"omitempty,max=64"`
	Price         int64          `json:"price" binding:"required,min=0"`
	OriginalPrice int64          `json:"original_price" binding:"omitempty,min=0"`
	Tags          string         `json:"tags" binding:"omitempty,max=256"`
	Images        string         `json:"images"`
	Description   string         `json:"description"`
	Origin        string         `json:"origin" binding:"omitempty,max=128"`
	ShelfLife     string         `json:"shelf_life" binding:"omitempty,max=64"`
	IsNew         bool           `json:"is_new"`
	IsHot         bool           `json:"is_hot"`
	IsFeatured    bool           `json:"is_featured"`
	SKUs          []CreateSKUReq `json:"skus"`
}

type CreateSKUReq struct {
	Name  string `json:"name" binding:"required,max=64"`
	Price int64  `json:"price" binding:"required,min=0"`
	Stock int    `json:"stock" binding:"min=0"`
}

type UpdateProductRequest struct {
	CategoryID    uint   `json:"category_id"`
	Name          string `json:"name" binding:"omitempty,max=128"`
	Subtitle      string `json:"subtitle" binding:"omitempty,max=256"`
	Brand         string `json:"brand" binding:"omitempty,max=64"`
	Price         *int64 `json:"price"`
	OriginalPrice *int64 `json:"original_price"`
	Tags          string `json:"tags" binding:"omitempty,max=256"`
	Images        string `json:"images"`
	Description   string `json:"description"`
	Origin        string `json:"origin" binding:"omitempty,max=128"`
	ShelfLife     string `json:"shelf_life" binding:"omitempty,max=64"`
	Status        string `json:"status" binding:"omitempty,oneof=on off"`
	IsNew         *bool  `json:"is_new"`
	IsHot         *bool  `json:"is_hot"`
	IsFeatured    *bool  `json:"is_featured"`
}

type ListProductRequest struct {
	CategoryID uint   `form:"category_id"`
	Tag        string `form:"tag"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort" binding:"omitempty,oneof=sales price_newest price_oldest newest"`
	Status     string `form:"status" binding:"omitempty,oneof=on off"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=50"`
}

type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Score     int    `json:"score" binding:"required,min=1,max=5"`
	Content   string `json:"content" binding:"omitempty,max=500"`
	Images    string `json:"images"`
}

type ListReviewRequest struct {
	ProductID uint `form:"product_id" binding:"required"`
	Page      int  `form:"page" binding:"omitempty,min=1"`
	Size      int  `form:"size" binding:"omitempty,min=1,max=50"`
}
