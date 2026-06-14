package browse_history

import (
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/browse_history"
)

type Service interface {
	Add(userID, productID uint) (*response.BrowseHistoryResponse, error)
	List(userID uint, req *request.ListBrowseHistoryRequest) ([]response.BrowseHistoryResponse, int64, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Add(userID, productID uint) (*response.BrowseHistoryResponse, error) {
	history := &model.BrowseHistory{
		UserID:    userID,
		ProductID: productID,
	}
	if err := s.repo.Create(history); err != nil {
		return nil, err
	}
	resp := response.ToBrowseHistoryResponse(history)
	return &resp, nil
}

func (s *service) List(userID uint, req *request.ListBrowseHistoryRequest) ([]response.BrowseHistoryResponse, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 || size > 50 {
		size = 10
	}
	list, total, err := s.repo.FindByUserID(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToBrowseHistoryList(list), total, nil
}
