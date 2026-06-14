package response

import (
	"gosh/internal/model"
)

type CartItemResponse struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	SKUID    uint   `json:"sku_id"`
	Quantity int    `json:"quantity"`
	Selected bool   `json:"selected"`
	// SKU 快照信息（查询时填充）
	ProductName string `json:"product_name"`
	SKUName     string `json:"sku_name"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	ProductID   uint   `json:"product_id"`
	ProductOn   bool   `json:"product_on"`
	CreatedAt   string `json:"created_at"`
}

type CartSummaryResponse struct {
	Items        []CartItemResponse `json:"items"`
	SelectedAll  bool               `json:"selected_all"`
	SelectedIDs  []uint             `json:"selected_ids"`
	TotalAmount  int64              `json:"total_amount"`
	TotalCount   int                `json:"total_count"`
	RemovedItems []RemovedItem      `json:"removed_items,omitempty"`
}

type RemovedItem struct {
	SKUID  uint   `json:"sku_id"`
	Reason string `json:"reason"`
}

func ToCartItemResponse(cart *model.Cart) CartItemResponse {
	return CartItemResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SKUID:     cart.SKUID,
		Quantity:  cart.Quantity,
		Selected:  cart.Selected,
		CreatedAt: cart.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToCartItemList(carts []model.Cart) []CartItemResponse {
	list := make([]CartItemResponse, len(carts))
	for i, c := range carts {
		list[i] = ToCartItemResponse(&c)
	}
	return list
}
