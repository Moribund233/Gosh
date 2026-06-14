package product

import (
	"errors"

	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/product"
	reviewRepo "gosh/internal/repository/review"
	searchRepo "gosh/internal/repository/search"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type Service interface {
	Create(req *request.CreateProductRequest) (*response.ProductResponse, error)
	Update(id uint, req *request.UpdateProductRequest) (*response.ProductResponse, error)
	GetByID(id uint) (*response.ProductResponse, error)
	List(req *request.ListProductRequest) ([]response.ProductResponse, int64, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
	Search(req *request.ListProductRequest, userID uint) ([]response.ProductResponse, int64, error)
	HotSearch() ([]model.HotSearch, error)
	ClearSearchHistory(userID uint) error
	SearchHistory(userID uint) ([]string, error)
}

type service struct {
	repo        repo.Repository
	reviewRepo  reviewRepo.Repository
	searchRepo  searchRepo.Repository
}

func New() Service {
	return &service{
		repo:       repo.New(),
		reviewRepo: reviewRepo.New(),
		searchRepo: searchRepo.New(),
	}
}

func (s *service) Create(req *request.CreateProductRequest) (*response.ProductResponse, error) {
	product := &model.Product{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		Brand:         req.Brand,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Tags:          req.Tags,
		Images:        req.Images,
		Description:   req.Description,
		Origin:        req.Origin,
		ShelfLife:     req.ShelfLife,
		IsNew:         req.IsNew,
		IsHot:         req.IsHot,
		IsFeatured:    req.IsFeatured,
		Status:        model.ProductStatusOn,
	}
	var skus []model.ProductSKU
	for _, s := range req.SKUs {
		skus = append(skus, model.ProductSKU{
			Name:  s.Name,
			Price: s.Price,
			Stock: s.Stock,
		})
	}
	if err := s.repo.Create(product, skus); err != nil {
		return nil, err
	}
	resp := response.ToProductResponse(product)
	resp.SKUs = response.ToSKUList(skus)
	return &resp, nil
}

func (s *service) Update(id uint, req *request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Subtitle != "" {
		product.Subtitle = req.Subtitle
	}
	if req.Brand != "" {
		product.Brand = req.Brand
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.OriginalPrice != nil {
		product.OriginalPrice = *req.OriginalPrice
	}
	if req.Tags != "" {
		product.Tags = req.Tags
	}
	if req.Images != "" {
		product.Images = req.Images
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Origin != "" {
		product.Origin = req.Origin
	}
	if req.ShelfLife != "" {
		product.ShelfLife = req.ShelfLife
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if req.IsNew != nil {
		product.IsNew = *req.IsNew
	}
	if req.IsHot != nil {
		product.IsHot = *req.IsHot
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	skus, _ := s.repo.FindSKUsByProductID(product.ID)
	resp := response.ToProductResponse(product)
	resp.SKUs = response.ToSKUList(skus)
	return &resp, nil
}

func (s *service) GetByID(id uint) (*response.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	resp := response.ToProductResponse(product)
	resp.SKUs = response.ToSKUList(product.SKUs)
	return &resp, nil
}

func (s *service) List(req *request.ListProductRequest) ([]response.ProductResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 || size > 50 {
		size = 10
	}
	products, total, err := s.repo.List(req.CategoryID, req.Tag, "", req.Sort, req.Status, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToProductList(products), total, nil
}

func (s *service) Search(req *request.ListProductRequest, userID uint) ([]response.ProductResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 || size > 50 {
		size = 10
	}
	products, total, err := s.repo.List(req.CategoryID, "", req.Keyword, req.Sort, "", page, size)
	if err != nil {
		return nil, 0, err
	}

	if userID > 0 && req.Keyword != "" {
		s.searchRepo.CreateHistory(&model.SearchHistory{UserID: userID, Query: req.Keyword})
		s.searchRepo.IncrementHotSearch(req.Keyword)
	}

	return response.ToProductList(products), total, nil
}

func (s *service) HotSearch() ([]model.HotSearch, error) {
	return s.searchRepo.HotSearch(10)
}

func (s *service) SearchHistory(userID uint) ([]string, error) {
	histories, err := s.searchRepo.FindHistoryByUserID(userID, 10)
	if err != nil {
		return nil, err
	}
	var queries []string
	for _, h := range histories {
		queries = append(queries, h.Query)
	}
	return queries, nil
}

func (s *service) ClearSearchHistory(userID uint) error {
	return s.searchRepo.DeleteHistoryByUserID(userID)
}

func (s *service) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrProductNotFound
	}
	return s.repo.Delete(id)
}

func (s *service) UpdateStatus(id uint, status string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrProductNotFound
	}
	return s.repo.UpdateStatus(id, status)
}
