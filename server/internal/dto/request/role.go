package request

type ListRoleUsersRequest struct {
	Role string `form:"role" binding:"omitempty"`
	Page int    `form:"page" binding:"omitempty,min=1"`
	Size int    `form:"size" binding:"omitempty,min=1,max=50"`
}

type UpdateUserRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=user merchant support operator super_admin"`
}
