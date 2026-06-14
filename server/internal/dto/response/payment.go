package response

import (
	"gosh/internal/model"
)

type PaymentMethodResponse struct {
	Method string `json:"method"`
	Name   string `json:"name"`
}

type PaymentResponse struct {
	ID            uint    `json:"id"`
	OrderNo       string  `json:"order_no"`
	TransactionNo string  `json:"transaction_no"`
	Method        string  `json:"method"`
	PayAmount     int64   `json:"pay_amount"`
	Status        string  `json:"status"`
	PaidAt        *string `json:"paid_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

func ToPaymentResponse(p *model.Payment) PaymentResponse {
	r := PaymentResponse{
		ID:            p.ID,
		OrderNo:       p.OrderNo,
		TransactionNo: p.TransactionNo,
		Method:        p.Method,
		PayAmount:     p.PayAmount,
		Status:        p.Status,
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.PaidAt != nil {
		s := p.PaidAt.Format("2006-01-02 15:04:05")
		r.PaidAt = &s
	}
	return r
}
