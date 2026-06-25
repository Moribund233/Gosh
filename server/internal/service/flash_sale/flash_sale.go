package flash_sale

import (
	"context"
	"time"

	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/flash_sale"
	"gosh/pkg/cache"
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
	c := cache.Default()
	if c != nil {
		var result []response.FlashSaleResponse
		err := c.Remember(context.Background(), "cache:flash_sale:active", 30*time.Second, func() (interface{}, error) {
			return s.buildActiveList()
		}, &result)
		if err == nil {
			return result, nil
		}
	}

	return s.buildActiveList()
}

func (s *service) buildActiveList() ([]response.FlashSaleResponse, error) {
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
