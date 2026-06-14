package request

type PayRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
	Method  string `json:"method" binding:"required,oneof=mock wechat alipay"`
}

type RefundRequest struct {
	TransactionNo string `json:"transaction_no" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=0"`
	Reason        string `json:"reason" binding:"omitempty,max=256"`
}
