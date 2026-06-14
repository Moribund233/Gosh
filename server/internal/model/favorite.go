package model

type Favorite struct {
	BaseModel
	UserID    uint  `gorm:"index:idx_user_product,not null" json:"user_id"`
	ProductID uint  `gorm:"index:idx_user_product,not null" json:"product_id"`
}
