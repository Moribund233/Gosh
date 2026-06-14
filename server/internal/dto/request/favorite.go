package request

type AddFavoriteRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type ListFavoriteRequest struct {
	Page int `form:"page" binding:"omitempty,min=1"`
	Size int `form:"size" binding:"omitempty,min=1,max=50"`
}
