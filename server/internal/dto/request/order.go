package request

type CreateOrderItem struct {
	SKUID    uint `json:"sku_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	AddressID      uint              `json:"address_id"`
	Remark         string            `json:"remark" binding:"omitempty,max=200"`
	DeliveryMethod string            `json:"delivery_method" binding:"omitempty,oneof=standard express"`
	Items          []CreateOrderItem `json:"items,omitempty"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=128"`
}

type ListOrderRequest struct {
	Status string `form:"status" binding:"omitempty,oneof=pending_payment pending_delivery pending_receipt completed cancelled"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=50"`
}

type ApplyPointsRequest struct {
	Points int `json:"points" binding:"required,min=1"`
}
