package merchant

import (
	"errors"

	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/merchant"
	userRepo "gosh/internal/repository/user"
)

var (
	ErrApplicationExists   = errors.New("application already exists")
	ErrApplicationNotFound = errors.New("application not found")
	ErrAlreadyReviewed     = errors.New("application already reviewed")
)

type Service interface {
	Apply(userID uint, req *request.ApplyMerchantRequest) (*response.MerchantApplicationResponse, error)
	Review(adminID uint, req *request.ReviewMerchantRequest) (*response.MerchantApplicationResponse, error)
	List(status string, page, size int) ([]response.MerchantApplicationResponse, int64, error)
	MyApplication(userID uint) (*response.MerchantApplicationResponse, error)
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

func (s *service) Apply(userID uint, req *request.ApplyMerchantRequest) (*response.MerchantApplicationResponse, error) {
	existing, _ := s.repo.FindByUserID(userID)
	if existing != nil && existing.Status == model.AppStatusPending {
		return nil, ErrApplicationExists
	}
	app := &model.MerchantApplication{
		UserID:       userID,
		ShopName:     req.ShopName,
		ShopDesc:     req.ShopDesc,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       model.AppStatusPending,
	}
	if err := s.repo.Create(app); err != nil {
		return nil, err
	}
	resp := response.ToMerchantApplicationResponse(app)
	return &resp, nil
}

func (s *service) Review(adminID uint, req *request.ReviewMerchantRequest) (*response.MerchantApplicationResponse, error) {
	app, err := s.repo.FindByID(req.ApplicationID)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	if app.Status != model.AppStatusPending {
		return nil, ErrAlreadyReviewed
	}
	switch req.Action {
	case "approve":
		app.Status = model.AppStatusApproved
		user, uErr := s.userRepo.FindByID(app.UserID)
		if uErr != nil {
			return nil, uErr
		}
		user.Role = model.RoleMerchant
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	case "reject":
		app.Status = model.AppStatusRejected
	}
	app.Remark = req.Remark
	if err := s.repo.Update(app); err != nil {
		return nil, err
	}
	resp := response.ToMerchantApplicationResponse(app)
	return &resp, nil
}

func (s *service) List(status string, page, size int) ([]response.MerchantApplicationResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	apps, total, err := s.repo.List(status, page, size)
	if err != nil {
		return nil, 0, err
	}
	return response.ToMerchantApplicationList(apps), total, nil
}

func (s *service) MyApplication(userID uint) (*response.MerchantApplicationResponse, error) {
	app, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	resp := response.ToMerchantApplicationResponse(app)
	return &resp, nil
}
