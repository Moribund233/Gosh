package model

type BrowseHistory struct {
	BaseModel
	UserID    uint `gorm:"index;not null" json:"user_id"`
	ProductID uint `gorm:"index;not null" json:"product_id"`
}
