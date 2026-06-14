package brand_story

import (
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/brand_story"
)

type Service interface {
	Get() (*response.BrandStoryResponse, error)
	Update(req *request.UpdateBrandStoryRequest) (*response.BrandStoryResponse, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Get() (*response.BrandStoryResponse, error) {
	story, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	resp := response.ToBrandStoryResponse(story)
	return &resp, nil
}

func (s *service) Update(req *request.UpdateBrandStoryRequest) (*response.BrandStoryResponse, error) {
	story := &model.BrandStory{
		Title:       req.Title,
		Description: req.Description,
		Logo:        req.Logo,
		Badge:       req.Badge,
		Link:        req.Link,
		Status:      model.StatusOn,
	}
	if err := s.repo.Upsert(story); err != nil {
		return nil, err
	}
	resp := response.ToBrandStoryResponse(story)
	return &resp, nil
}
