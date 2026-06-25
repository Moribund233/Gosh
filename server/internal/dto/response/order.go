package response

import (
	"encoding/json"
	"time"

	"gosh/internal/model"
)

type OrderResponse struct {
	ID             uint                `json:"id"`
	OrderNo        string              `json:"order_no"`
	UserID         uint                `json:"user_id"`
	Status         string              `json:"status"`
	TotalAmount    int64               `json:"total_amount"`
	ShippingFee    int64               `json:"shipping_fee"`
	DiscountAmount int64               `json:"discount_amount"`
	PayAmount      int64               `json:"pay_amount"`
	Remark         string              `json:"remark"`
	DeliveryMethod string              `json:"delivery_method"`
	AddressName    string              `json:"address_name"`
	AddressPhone   string              `json:"address_phone"`
	AddressDetail  string              `json:"address_detail"`
	PaidAt         *string             `json:"paid_at,omitempty"`
	ShippedAt      *string             `json:"shipped_at,omitempty"`
	CompletedAt    *string             `json:"completed_at,omitempty"`
	CancelledAt    *string             `json:"cancelled_at,omitempty"`
	CancelReason   string              `json:"cancel_reason,omitempty"`
	CreatedAt      string              `json:"created_at"`
	Items          []OrderItemResponse `json:"items,omitempty"`
}

type OrderItemResponse struct {
	ID          uint   `json:"id"`
	OrderID     uint   `json:"order_id"`
	SKUID       uint   `json:"sku_id"`
	ProductName string `json:"product_name"`
	SKUName     string `json:"sku_name"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}

type RebuyResponse struct {
	Cart         CartListResponse `json:"cart"`
	SkippedItems []SkippedItem    `json:"skipped_items"`
}

type SkippedItem struct {
	SKUID  uint   `json:"sku_id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type CartListResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCount int                `json:"total_count"`
}

type AddressSnapshot struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Province string `json:"province"`
	City    string `json:"city"`
	District string `json:"district"`
	Detail  string `json:"detail"`
}

func ToOrderResponse(o *model.Order) OrderResponse {
	r := OrderResponse{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		UserID:         o.UserID,
		Status:         o.Status,
		TotalAmount:    o.TotalAmount,
		ShippingFee:    o.ShippingFee,
		DiscountAmount: o.DiscountAmount,
		PayAmount:      o.PayAmount,
		Remark:         o.Remark,
		DeliveryMethod: o.DeliveryMethod,
		AddressName:    o.AddressName,
		AddressPhone:   o.AddressPhone,
		AddressDetail:  o.AddressDetail,
		CancelReason:   o.CancelReason,
		CreatedAt:      o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if o.PaidAt != nil {
		s := o.PaidAt.Format("2006-01-02 15:04:05")
		r.PaidAt = &s
	}
	if o.ShippedAt != nil {
		s := o.ShippedAt.Format("2006-01-02 15:04:05")
		r.ShippedAt = &s
	}
	if o.CompletedAt != nil {
		s := o.CompletedAt.Format("2006-01-02 15:04:05")
		r.CompletedAt = &s
	}
	if o.CancelledAt != nil {
		s := o.CancelledAt.Format("2006-01-02 15:04:05")
		r.CancelledAt = &s
	}
	if len(o.Items) > 0 {
		r.Items = make([]OrderItemResponse, len(o.Items))
		for i, item := range o.Items {
			r.Items[i] = ToOrderItemResponse(&item)
		}
	}
	return r
}

func ToOrderItemResponse(item *model.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:          item.ID,
		OrderID:     item.OrderID,
		SKUID:       item.SKUID,
		ProductName: item.ProductName,
		SKUName:     item.SKUName,
		Image:       item.Image,
		Price:       item.Price,
		Quantity:    item.Quantity,
		Subtotal:    item.Subtotal,
	}
}

func ToOrderList(orders []model.Order) []OrderResponse {
	list := make([]OrderResponse, len(orders))
	for i, o := range orders {
		list[i] = ToOrderResponse(&o)
	}
	return list
}

func AddressToDetail(addr *model.Address) string {
	snap := AddressSnapshot{
		Name:    addr.Name,
		Phone:   addr.Phone,
		Province: addr.Province,
		City:    addr.City,
		District: addr.District,
		Detail:  addr.Detail,
	}
	b, _ := json.Marshal(snap)
	return string(b)
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
