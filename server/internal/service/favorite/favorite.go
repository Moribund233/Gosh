package favorite

import (
	"errors"

	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/favorite"
)

var (
	ErrFavoriteExists = errors.New("already favorited")
	ErrFavoriteNotFound = errors.New("favorite not found")
)

type Service interface {
	Add(userID, productID uint) (*response.FavoriteResponse, error)
	Remove(userID, productID uint) error
	List(userID uint, req *request.ListFavoriteRequest) ([]response.FavoriteResponse, int64, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Add(userID, productID uint) (*response.FavoriteResponse, error) {
	exists, _ := s.repo.Exists(userID, productID)
	if exists {
		return nil, ErrFavoriteExists
	}
	fav := &model.Favorite{
		UserID:    userID,
		ProductID: productID,
	}
	if err := s.repo.Create(fav); err != nil {
		return nil, err
	}
	resp := response.ToFavoriteResponse(fav)
	return &resp, nil
}

func (s *service) Remove(userID, productID uint) error {
	exists, _ := s.repo.Exists(userID, productID)
	if !exists {
		return ErrFavoriteNotFound
	}
	return s.repo.Delete(userID, productID)
}

func (s *service) List(userID uint, req *request.ListFavoriteRequest) ([]response.FavoriteResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 || size > 50 {
		size = 10
	}
	favs, total, err := s.repo.FindByUserID(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToFavoriteList(favs), total, nil
}
