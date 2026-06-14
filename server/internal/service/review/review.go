package review

import (
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/review"
	userRepo "gosh/internal/repository/user"
)

type Service interface {
	Create(userID uint, req *request.CreateReviewRequest) (*response.ReviewResponse, error)
	List(productID uint, page, size int) ([]response.ReviewResponse, int64, error)
}

type service struct {
	repo     repo.Repository
	userRepo userRepo.Repository
}

func New() Service {
	return &service{
		repo:     repo.New(),
		userRepo: userRepo.New(),
	}
}

func (s *service) Create(userID uint, req *request.CreateReviewRequest) (*response.ReviewResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	nickname := user.Nickname
	if len([]rune(nickname)) > 1 {
		runes := []rune(nickname)
		nickname = string(runes[0]) + "***"
	}
	review := &model.ProductReview{
		ProductID: req.ProductID,
		UserID:    userID,
		Score:     req.Score,
		Content:   req.Content,
		Images:    req.Images,
		Nickname:  nickname,
		Avatar:    user.Avatar,
	}
	if err := s.repo.Create(review); err != nil {
		return nil, err
	}
	resp := response.ToReviewResponse(review)
	return &resp, nil
}

func (s *service) List(productID uint, page, size int) ([]response.ReviewResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	reviews, total, err := s.repo.FindByProductID(productID, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToReviewList(reviews), total, nil
}
