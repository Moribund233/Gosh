package flash_sale

import (
	"time"

	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/flash_sale"
)

type Service interface {
	GetActive() ([]response.FlashSaleResponse, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) GetActive() ([]response.FlashSaleResponse, error) {
	list, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}

	resp := make([]response.FlashSaleResponse, len(list))
	now := time.Now()
	for i, fs := range list {
		r := response.ToFlashSaleResponse(&fs)
		r.Countdown = int64(fs.EndAt.Sub(now).Seconds())
		if r.Countdown < 0 {
			r.Countdown = 0
		}
		resp[i] = r
	}
	return resp, nil
}

func (s *service) GetActiveForProduct(productID uint) (*model.FlashSale, error) {
	return s.repo.FindByProductID(productID)
}
