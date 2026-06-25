package model

import "time"

type FlashSale struct {
	BaseModel
	ProductID  uint      `gorm:"index;not null" json:"product_id"`
	SKUID      uint      `gorm:"index;not null" json:"sku_id"`
	FlashPrice int64     `gorm:"not null" json:"flash_price"`
	FlashStock int       `gorm:"not null" json:"flash_stock"`
	StartAt    time.Time `gorm:"index:idx_flash_active;not null" json:"start_at"`
	EndAt      time.Time `gorm:"index:idx_flash_active;not null" json:"end_at"`
	Status     string    `gorm:"size:20;default:active;index:idx_flash_active" json:"status"`
	Version    int       `gorm:"default:0" json:"version"`
}

const (
	FlashSaleStatusActive   = "active"
	FlashSaleStatusExpired  = "expired"
	FlashSaleStatusDisabled = "disabled"
)
