package model

import "time"

const (
	PaymentMethodMock   = "mock"
	PaymentMethodWechat = "wechat"
	PaymentMethodAlipay = "alipay"

	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

type Payment struct {
	BaseModel
	OrderID       uint       `gorm:"index;not null" json:"order_id"`
	OrderNo       string     `gorm:"size:32;not null" json:"order_no"`
	Method        string     `gorm:"size:20;not null" json:"method"`
	PayAmount     int64      `gorm:"not null" json:"pay_amount"`
	Status        string     `gorm:"size:20;default:pending;index" json:"status"`
	TransactionNo string     `gorm:"size:64;uniqueIndex" json:"transaction_no"`
	PaidAt        *time.Time `json:"paid_at"`
	NotifyRaw     string     `gorm:"type:text" json:"notify_raw"`
	NotifySignOk  bool       `gorm:"default:false" json:"notify_sign_ok"`
}
