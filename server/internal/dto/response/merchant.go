package response

import (
	"gosh/internal/model"
)

type MerchantApplicationResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	ShopName     string `json:"shop_name"`
	ShopDesc     string `json:"shop_desc"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
}

func ToMerchantApplicationResponse(a *model.MerchantApplication) MerchantApplicationResponse {
	return MerchantApplicationResponse{
		ID:           a.ID,
		UserID:       a.UserID,
		ShopName:     a.ShopName,
		ShopDesc:     a.ShopDesc,
		ContactName:  a.ContactName,
		ContactPhone: a.ContactPhone,
		Status:       a.Status,
		Remark:       a.Remark,
		CreatedAt:    a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToMerchantApplicationList(apps []model.MerchantApplication) []MerchantApplicationResponse {
	list := make([]MerchantApplicationResponse, len(apps))
	for i, a := range apps {
		list[i] = ToMerchantApplicationResponse(&a)
	}
	return list
}
