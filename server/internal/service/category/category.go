package category

import (
	"errors"

	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/category"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrHasChildren      = errors.New("category has children, cannot delete")
)

type Service interface {
	Create(req *request.CreateCategoryRequest) (*response.CategoryResponse, error)
	Update(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error)
	Delete(id uint) error
	GetTree() ([]response.CategoryResponse, error)
	GetByID(id uint) (*response.CategoryResponse, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Create(req *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	cat := &model.Category{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Icon:      req.Icon,
		Banner:    req.Banner,
		SortOrder: req.SortOrder,
	}
	if req.ParentID != nil {
		parent, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			return nil, ErrCategoryNotFound
		}
		cat.Level = parent.Level + 1
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	resp := response.ToCategoryResponse(cat)
	return &resp, nil
}

func (s *service) Update(id uint, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Icon != "" {
		cat.Icon = req.Icon
	}
	if req.Banner != "" {
		cat.Banner = req.Banner
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	resp := response.ToCategoryResponse(cat)
	return &resp, nil
}

func (s *service) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrCategoryNotFound
	}
	count, err := s.repo.CountByParentID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrHasChildren
	}
	return s.repo.Delete(id)
}

func (s *service) GetTree() ([]response.CategoryResponse, error) {
	cats, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return response.ToCategoryTree(cats), nil
}

func (s *service) GetByID(id uint) (*response.CategoryResponse, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	resp := response.ToCategoryResponse(cat)
	return &resp, nil
}
