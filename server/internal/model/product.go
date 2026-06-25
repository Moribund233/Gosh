package model

type Product struct {
	BaseModel
	CategoryID    uint        `gorm:"index;not null" json:"category_id"`
	Name          string      `gorm:"size:128;not null" json:"name"`
	Subtitle      string      `gorm:"size:256" json:"subtitle"`
	Brand         string      `gorm:"size:64" json:"brand"`
	Price         int64       `gorm:"not null" json:"price"`
	OriginalPrice int64       `json:"original_price"`
	Sales         int64       `gorm:"default:0" json:"sales"`
	Tags          string      `gorm:"size:256" json:"tags"`
	Images        []string    `gorm:"type:text;serializer:json" json:"images"`
	Status        string      `gorm:"size:20;default:on;index" json:"status"`
	Description   string      `gorm:"type:text" json:"description"`
	Origin        string      `gorm:"size:128" json:"origin"`
	ShelfLife     string      `gorm:"size:64" json:"shelf_life"`
	IsNew         bool        `gorm:"default:false;index" json:"is_new"`
	IsHot         bool        `gorm:"default:false;index" json:"is_hot"`
	IsFeatured    bool        `gorm:"default:false;index" json:"is_featured"`
	SKUs          []ProductSKU `json:"skus,omitempty"`
}

type ProductSKU struct {
	BaseModel
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	Name      string `gorm:"size:64;not null" json:"name"`
	Price     int64  `gorm:"not null" json:"price"`
	Stock     int    `gorm:"default:0" json:"stock"`
	Version   int    `gorm:"default:0" json:"version"`
}

const (
	ProductStatusOn  = "on"
	ProductStatusOff = "off"
)

type ProductReview struct {
	BaseModel
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	OrderID   uint   `gorm:"index" json:"order_id"`
	Score     int    `gorm:"not null" json:"score"`
	Content   string `gorm:"type:text" json:"content"`
	Images    []string `gorm:"type:text;serializer:json" json:"images"`
	Nickname  string `gorm:"size:64" json:"nickname"`
	Avatar    string `gorm:"size:256" json:"avatar"`
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

type SearchHistory struct {
	BaseModel
	UserID uint   `gorm:"index;not null" json:"user_id"`
	Query  string `gorm:"size:128;not null" json:"query"`
}

func (SearchHistory) TableName() string {
	return "search_histories"
}

type HotSearch struct {
	BaseModel
	Query     string `gorm:"size:128;not null;uniqueIndex" json:"query"`
	Count     int64  `gorm:"default:0" json:"count"`
}
