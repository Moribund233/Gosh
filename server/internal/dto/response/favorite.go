package response

import (
	"gosh/internal/model"
)

type FavoriteResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	CreatedAt string `json:"created_at"`
}

func ToFavoriteResponse(f *model.Favorite) FavoriteResponse {
	return FavoriteResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		ProductID: f.ProductID,
		CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToFavoriteList(favs []model.Favorite) []FavoriteResponse {
	list := make([]FavoriteResponse, len(favs))
	for i, f := range favs {
		list[i] = ToFavoriteResponse(&f)
	}
	return list
}
