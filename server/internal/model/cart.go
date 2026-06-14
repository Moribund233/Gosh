package model

type Cart struct {
	BaseModel
	UserID   uint `gorm:"uniqueIndex:idx_cart_user_sku;not null;column:user_id" json:"user_id"`
	SKUID    uint `gorm:"uniqueIndex:idx_cart_user_sku;not null;column:sku_id" json:"sku_id"`
	Quantity int  `gorm:"not null;default:1;column:quantity" json:"quantity"`
	Selected bool `gorm:"default:true;column:selected" json:"selected"`
}
