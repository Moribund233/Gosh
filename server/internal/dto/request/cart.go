package request

type AddCartRequest struct {
	SKUID    uint `json:"sku_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartRequest struct {
	Quantity int  `json:"quantity" binding:"required,min=1"`
	Selected *bool `json:"selected"`
}

type SelectRequest struct {
	SKUIDs []uint `json:"sku_ids"` // 空列表=全选/全不选
	All    *bool  `json:"all"`     // true=全选 false=取消全选
	Select bool   `json:"select"`  // true=选中 false=取消选中
}

type MergeCartRequest struct {
	Items []MergeItem `json:"items" binding:"required"`
}

type MergeItem struct {
	SKUID    uint `json:"sku_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}
