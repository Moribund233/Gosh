package model

const (
	PointTypeEarn = "earn"
	PointTypeSpend = "spend"
)

type PointLog struct {
	BaseModel
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	Type    string `gorm:"size:20;not null" json:"type"`
	Amount  int    `gorm:"not null" json:"amount"`
	Balance int    `gorm:"not null" json:"balance"`
	OrderID *uint  `gorm:"index" json:"order_id"`
	Note    string `gorm:"size:256" json:"note"`
}

func (PointLog) TableName() string {
	return "point_logs"
}
