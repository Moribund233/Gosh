package model

type Address struct {
	BaseModel
	UserID    uint   `gorm:"index:idx_address_user_default;not null" json:"user_id"`
	Name      string `gorm:"size:32;not null" json:"name"`
	Phone     string `gorm:"size:20;not null" json:"phone"`
	Province  string `gorm:"size:32;not null" json:"province"`
	City      string `gorm:"size:32;not null" json:"city"`
	District  string `gorm:"size:32;not null" json:"district"`
	Detail    string `gorm:"size:256;not null" json:"detail"`
	IsDefault bool   `gorm:"default:false;index:idx_address_user_default" json:"is_default"`
}
