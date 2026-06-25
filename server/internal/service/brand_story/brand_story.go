package brand_story

import (
	"context"
	"time"

	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/brand_story"
	"gosh/pkg/cache"
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
	c := cache.Default()
	if c != nil {
		var result response.BrandStoryResponse
		err := c.Remember(context.Background(), "cache:brand_story", 1*time.Hour, func() (interface{}, error) {
			story, err := s.repo.Get()
			if err != nil {
				return nil, err
			}
			return response.ToBrandStoryResponse(story), nil
		}, &result)
		if err == nil {
			return &result, nil
		}
	}

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
	if c := cache.Default(); c != nil {
		c.Del(context.Background(), "cache:brand_story")
	}
	resp := response.ToBrandStoryResponse(story)
	return &resp, nil
}
