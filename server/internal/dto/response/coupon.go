package response

import (
	"gosh/internal/model"
)

type CouponResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Condition   int64  `json:"condition"`
	Discount    int64  `json:"discount"`
	TotalCount  int    `json:"total_count"`
	RemainCount int    `json:"remain_count"`
	PerLimit    int    `json:"per_limit"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	Status      string `json:"status"`
}

type UserCouponResponse struct {
	ID         uint   `json:"id"`
	CouponID   uint   `json:"coupon_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Condition  int64  `json:"condition"`
	Discount   int64  `json:"discount"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	Status     string `json:"status"`
	UsedAt     string `json:"used_at,omitempty"`
}

type CouponCalculateResponse struct {
	OriginalAmount int64 `json:"original_amount"`
	DiscountAmount int64 `json:"discount_amount"`
	PayAmount      int64 `json:"pay_amount"`
}

func ToCouponResponse(c *model.Coupon) CouponResponse {
	return CouponResponse{
		ID:          c.ID,
		Name:        c.Name,
		Type:        c.Type,
		Condition:   c.Condition,
		Discount:    c.Discount,
		TotalCount:  c.TotalCount,
		RemainCount: c.RemainCount,
		PerLimit:    c.PerLimit,
		StartAt:     c.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:       c.EndAt.Format("2006-01-02 15:04:05"),
		Status:      c.Status,
	}
}

func ToUserCouponResponse(uc *model.UserCoupon, c *model.Coupon) UserCouponResponse {
	r := UserCouponResponse{
		ID:        uc.ID,
		CouponID:  c.ID,
		Name:      c.Name,
		Type:      c.Type,
		Condition: c.Condition,
		Discount:  c.Discount,
		StartAt:   c.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:     c.EndAt.Format("2006-01-02 15:04:05"),
		Status:    uc.Status,
	}
	if uc.UsedAt != nil {
		r.UsedAt = uc.UsedAt.Format("2006-01-02 15:04:05")
	}
	return r
}
