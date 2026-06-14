package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/user"
	favRepo "gosh/internal/repository/favorite"
	histRepo "gosh/internal/repository/browse_history"
	"gosh/pkg/auth"
)

var (
	ErrPhoneExists  = errors.New("phone already registered")
	ErrInvalidCreds = errors.New("invalid phone or password")
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	Register(phone, password, nickname string) (*response.TokenResponse, error)
	Login(phone, password string) (*response.TokenResponse, error)
	GetProfile(userID uint) (*response.ProfileResponse, error)
	UpdateProfile(userID uint, req *request.UpdateProfileRequest) (*response.UserResponse, error)
	ListByRole(role string, page, size int) ([]response.UserResponse, int64, error)
	UpdateRole(userID uint, role string) error
}

type service struct {
	repo       repo.Repository
	favRepo    favRepo.Repository
	histRepo   histRepo.Repository
}

func New() Service {
	return &service{
		repo:     repo.New(),
		favRepo:  favRepo.New(),
		histRepo: histRepo.New(),
	}
}

func (s *service) Register(phone, password, nickname string) (*response.TokenResponse, error) {
	existing, _ := s.repo.FindByPhone(phone)
	if existing != nil {
		return nil, ErrPhoneExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:    phone,
		Password: string(hash),
		Nickname: nickname,
		Role:     model.RoleUser,
		Status:   model.StatusActive,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	token, expiresAt, err := auth.Sign(user.ID, user.Role, user.TenantID)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      response.ToUserResponse(user),
	}, nil
}

func (s *service) Login(phone, password string) (*response.TokenResponse, error) {
	user, err := s.repo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}

	if user.Status != model.StatusActive {
		return nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCreds
	}

	token, expiresAt, err := auth.Sign(user.ID, user.Role, user.TenantID)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      response.ToUserResponse(user),
	}, nil
}

func (s *service) GetProfile(userID uint) (*response.ProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	favCount, _ := s.favRepo.CountByUserID(userID)
	viewCount, _ := s.histRepo.CountByUserID(userID)
	return &response.ProfileResponse{
		User:      response.ToUserResponse(user),
		FavCount:  int(favCount),
		ViewCount: int(viewCount),
	}, nil
}

func (s *service) UpdateProfile(userID uint, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	resp := response.ToUserResponse(user)
	return &resp, nil
}

func (s *service) ListByRole(role string, page, size int) ([]response.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	users, total, err := s.repo.ListByRole(role, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.UserResponse, len(users))
	for i, u := range users {
		list[i] = response.ToUserResponse(&u)
	}
	return list, total, nil
}

func (s *service) UpdateRole(userID uint, role string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	user.Role = role
	return s.repo.Update(user)
}
