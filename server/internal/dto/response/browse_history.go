package response

import (
	"gosh/internal/model"
)

type BrowseHistoryResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	CreatedAt string `json:"created_at"`
}

func ToBrowseHistoryResponse(b *model.BrowseHistory) BrowseHistoryResponse {
	return BrowseHistoryResponse{
		ID:        b.ID,
		UserID:    b.UserID,
		ProductID: b.ProductID,
		CreatedAt: b.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToBrowseHistoryList(bh []model.BrowseHistory) []BrowseHistoryResponse {
	list := make([]BrowseHistoryResponse, len(bh))
	for i, b := range bh {
		list[i] = ToBrowseHistoryResponse(&b)
	}
	return list
}
