package request

type AddBrowseHistoryRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type ListBrowseHistoryRequest struct {
	Page int `form:"page" binding:"omitempty,min=1"`
	Size int `form:"size" binding:"omitempty,min=1,max=50"`
}
