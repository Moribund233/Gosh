package request

type CreateAddressRequest struct {
	Name      string `json:"name" binding:"required,max=32"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Province  string `json:"province" binding:"required,max=32"`
	City      string `json:"city" binding:"required,max=32"`
	District  string `json:"district" binding:"required,max=32"`
	Detail    string `json:"detail" binding:"required,max=256"`
	IsDefault bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Name      string `json:"name" binding:"omitempty,max=32"`
	Phone     string `json:"phone" binding:"omitempty,len=11"`
	Province  string `json:"province" binding:"omitempty,max=32"`
	City      string `json:"city" binding:"omitempty,max=32"`
	District  string `json:"district" binding:"omitempty,max=32"`
	Detail    string `json:"detail" binding:"omitempty,max=256"`
	IsDefault *bool  `json:"is_default"`
}
