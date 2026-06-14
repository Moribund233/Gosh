package banner

import (
	"errors"

	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/banner"
)

var (
	ErrBannerNotFound = errors.New("banner not found")
)

type Service interface {
	Create(req *request.CreateBannerRequest) (*response.BannerResponse, error)
	Update(id uint, req *request.UpdateBannerRequest) (*response.BannerResponse, error)
	Delete(id uint) error
	List(status string) ([]response.BannerResponse, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Create(req *request.CreateBannerRequest) (*response.BannerResponse, error) {
	banner := &model.Banner{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Image:       req.Image,
		Background:  req.Background,
		Link:        req.Link,
		SortOrder:   req.SortOrder,
		Status:      model.StatusOn,
	}
	if err := s.repo.Create(banner); err != nil {
		return nil, err
	}
	resp := response.ToBannerResponse(banner)
	return &resp, nil
}

func (s *service) Update(id uint, req *request.UpdateBannerRequest) (*response.BannerResponse, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrBannerNotFound
	}
	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.Subtitle != "" {
		banner.Subtitle = req.Subtitle
	}
	if req.Description != "" {
		banner.Description = req.Description
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Background != "" {
		banner.Background = req.Background
	}
	if req.Link != "" {
		banner.Link = req.Link
	}
	if req.SortOrder != nil {
		banner.SortOrder = *req.SortOrder
	}
	if req.Status != "" {
		banner.Status = req.Status
	}
	if err := s.repo.Update(banner); err != nil {
		return nil, err
	}
	resp := response.ToBannerResponse(banner)
	return &resp, nil
}

func (s *service) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrBannerNotFound
	}
	return s.repo.Delete(id)
}

func (s *service) List(status string) ([]response.BannerResponse, error) {
	list, err := s.repo.List(status)
	if err != nil {
		return nil, err
	}
	return response.ToBannerList(list), nil
}
