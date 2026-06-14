package request

type CreateCategoryRequest struct {
	ParentID  *uint  `json:"parent_id"`
	Name      string `json:"name" binding:"required,max=32"`
	Icon      string `json:"icon" binding:"omitempty,max=32"`
	Banner    string `json:"banner" binding:"omitempty,max=256"`
	SortOrder int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name" binding:"omitempty,max=32"`
	Icon      string `json:"icon" binding:"omitempty,max=32"`
	Banner    string `json:"banner" binding:"omitempty,max=256"`
	SortOrder *int   `json:"sort_order"`
}
