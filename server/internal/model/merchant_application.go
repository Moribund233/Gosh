package model

type MerchantApplication struct {
	BaseModel
	UserID       uint   `gorm:"index;not null" json:"user_id"`
	ShopName     string `gorm:"size:64;not null" json:"shop_name"`
	ShopDesc     string `gorm:"size:512" json:"shop_desc"`
	ContactName  string `gorm:"size:32;not null" json:"contact_name"`
	ContactPhone string `gorm:"size:20;not null" json:"contact_phone"`
	Status       string `gorm:"size:20;default:pending;index" json:"status"`
	Remark       string `gorm:"size:256" json:"remark"`
}

const (
	AppStatusPending  = "pending"
	AppStatusApproved = "approved"
	AppStatusRejected = "rejected"
)
