package model

import "time"

const (
	CouponTypeFullReduce = "full_reduce"
	CouponTypeDiscount   = "discount"

	CouponStatusActive   = "active"
	CouponStatusExpired  = "expired"
	CouponStatusDisabled = "disabled"

	UserCouponStatusUnused  = "unused"
	UserCouponStatusUsed    = "used"
	UserCouponStatusExpired = "expired"
)

type Coupon struct {
	BaseModel
	Name        string    `gorm:"size:64;not null" json:"name"`
	Type        string    `gorm:"size:20;not null" json:"type"`
	Condition   int64     `gorm:"not null" json:"condition"`
	Discount    int64     `gorm:"not null" json:"discount"`
	TotalCount  int       `gorm:"default:0" json:"total_count"`
	RemainCount int       `gorm:"default:0" json:"remain_count"`
	PerLimit    int       `gorm:"default:1" json:"per_limit"`
	StartAt     time.Time `gorm:"not null" json:"start_at"`
	EndAt       time.Time `gorm:"not null" json:"end_at"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
}

type UserCoupon struct {
	BaseModel
	UserID   uint    `gorm:"index;not null" json:"user_id"`
	CouponID uint    `gorm:"index;not null" json:"coupon_id"`
	UsedAt   *time.Time `json:"used_at"`
	OrderID  *uint   `gorm:"index" json:"order_id"`
	Status   string  `gorm:"size:20;default:unused" json:"status"`
}
