package model

import "time"

const (
	OrderStatusPendingPayment  = "pending_payment"
	OrderStatusPendingDelivery = "pending_delivery"
	OrderStatusPendingReceipt  = "pending_receipt"
	OrderStatusCompleted       = "completed"
	OrderStatusCancelled       = "cancelled"
)

type Order struct {
	BaseModel
	OrderNo        string     `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	Status         string     `gorm:"size:20;default:pending_payment;index" json:"status"`
	TotalAmount    int64      `gorm:"not null" json:"total_amount"`
	ShippingFee    int64      `gorm:"default:0" json:"shipping_fee"`
	DiscountAmount int64      `gorm:"default:0" json:"discount_amount"`
	PayAmount      int64      `gorm:"not null" json:"pay_amount"`
	Remark         string     `gorm:"size:200" json:"remark"`
	DeliveryMethod string     `gorm:"size:32" json:"delivery_method"`
	AddressName    string     `gorm:"size:32" json:"address_name"`
	AddressPhone   string     `gorm:"size:20" json:"address_phone"`
	AddressDetail  string     `gorm:"size:512" json:"address_detail"`
	PaidAt         *time.Time `json:"paid_at"`
	ShippedAt      *time.Time `json:"shipped_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	CancelledAt    *time.Time `json:"cancelled_at"`
	CancelReason   string     `gorm:"size:128" json:"cancel_reason"`
	PointsDeducted int        `gorm:"default:0" json:"points_deducted"`
	CreatedAt      time.Time  `gorm:"index" json:"created_at"`
	Version        int        `gorm:"default:0" json:"version"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	BaseModel
	OrderID     uint   `gorm:"index;not null" json:"order_id"`
	SKUID       uint   `gorm:"index;not null" json:"sku_id"`
	ProductName string `gorm:"size:128;not null" json:"product_name"`
	SKUName     string `gorm:"size:64" json:"sku_name"`
	Image       string `gorm:"size:256" json:"image"`
	Price       int64  `gorm:"not null" json:"price"`
	Quantity    int    `gorm:"not null" json:"quantity"`
	Subtotal    int64  `gorm:"not null" json:"subtotal"`
}

type OrderLog struct {
	BaseModel
	OrderID    uint   `gorm:"index:idx_order_logs;not null" json:"order_id"`
	FromStatus string `gorm:"size:20" json:"from_status"`
	ToStatus   string `gorm:"size:20;not null" json:"to_status"`
	Operator   string `gorm:"size:32" json:"operator"`
	Note       string `gorm:"size:256" json:"note"`
}

type IdempotencyRecord struct {
	BaseModel
	Key      string `gorm:"uniqueIndex;size:64;not null" json:"key"`
	Response string `gorm:"type:text" json:"response"`
}
