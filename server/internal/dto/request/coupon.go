package request

type CreateCouponRequest struct {
	Name       string `json:"name" binding:"required,max=64"`
	Type       string `json:"type" binding:"required,oneof=full_reduce discount"`
	Condition  int64  `json:"condition" binding:"required,min=0"`
	Discount   int64  `json:"discount" binding:"required,min=1"`
	TotalCount int    `json:"total_count" binding:"min=0"`
	PerLimit   int    `json:"per_limit" binding:"min=1"`
	StartAt    string `json:"start_at" binding:"required"`
	EndAt      string `json:"end_at" binding:"required"`
}

type CalculateCouponRequest struct {
	OrderAmount int64 `json:"order_amount" binding:"required,min=0"`
	CouponID    uint  `json:"coupon_id" binding:"required"`
}
