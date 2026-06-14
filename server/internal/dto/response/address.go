package response

import (
	"gosh/internal/model"
)

type AddressResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Detail     string `json:"detail"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
}

func ToAddressResponse(a *model.Address) AddressResponse {
	return AddressResponse{
		ID:        a.ID,
		Name:      a.Name,
		Phone:     a.Phone,
		Province:  a.Province,
		City:      a.City,
		District:  a.District,
		Detail:    a.Detail,
		IsDefault: a.IsDefault,
		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToAddressList(addrs []model.Address) []AddressResponse {
	list := make([]AddressResponse, len(addrs))
	for i, a := range addrs {
		list[i] = ToAddressResponse(&a)
	}
	return list
}
