package response

import (
	"gosh/internal/model"
)

type FlashSaleResponse struct {
	ID         uint   `json:"id"`
	ProductID  uint   `json:"product_id"`
	SKUID      uint   `json:"sku_id"`
	FlashPrice int64  `json:"flash_price"`
	FlashStock int    `json:"flash_stock"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	Countdown  int64  `json:"countdown"`
}

func ToFlashSaleResponse(fs *model.FlashSale) FlashSaleResponse {
	return FlashSaleResponse{
		ID:         fs.ID,
		ProductID:  fs.ProductID,
		SKUID:      fs.SKUID,
		FlashPrice: fs.FlashPrice,
		FlashStock: fs.FlashStock,
		StartAt:    fs.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:      fs.EndAt.Format("2006-01-02 15:04:05"),
		Countdown:  int64(fs.EndAt.Sub(fs.StartAt).Seconds()),
	}
}
