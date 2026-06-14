package response

import (
	"time"
	"gosh/internal/model"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Phone     string    `json:"phone"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`

	Status string `json:"status,omitempty"`
}

type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type ProfileResponse struct {
	User       UserResponse      `json:"user"`
	FavCount   int               `json:"fav_count"`
	ViewCount  int               `json:"view_count"`
	OrderStats OrderStatsResponse `json:"order_stats"`
}

type OrderStatsResponse struct {
	Unpaid      int `json:"unpaid"`
	Undelivered int `json:"undelivered"`
	Delivering  int `json:"delivering"`
	Completed   int `json:"completed"`
	AfterSale   int `json:"after_sale"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Phone:     u.Phone,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Role:      u.Role,
		Points:    u.Points,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
