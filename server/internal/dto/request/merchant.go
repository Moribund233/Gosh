package request

type ApplyMerchantRequest struct {
	ShopName     string `json:"shop_name" binding:"required,max=64"`
	ShopDesc     string `json:"shop_desc" binding:"omitempty,max=512"`
	ContactName  string `json:"contact_name" binding:"required,max=32"`
	ContactPhone string `json:"contact_phone" binding:"required,len=11"`
}

type ReviewMerchantRequest struct {
	ApplicationID uint   `json:"application_id" binding:"required"`
	Action        string `json:"action" binding:"required,oneof=approve reject"`
	Remark        string `json:"remark" binding:"omitempty,max=256"`
}
