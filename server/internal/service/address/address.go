package address

import (
	"errors"

	"gorm.io/gorm"
	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/address"
)

var (
	ErrAddressNotFound = errors.New("address not found")
)

type Service interface {
	Create(userID uint, req *request.CreateAddressRequest) (*response.AddressResponse, error)
	List(userID uint) ([]response.AddressResponse, error)
	Update(userID, addressID uint, req *request.UpdateAddressRequest) (*response.AddressResponse, error)
	Delete(userID, addressID uint) error
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Create(userID uint, req *request.CreateAddressRequest) (*response.AddressResponse, error) {
	if req.IsDefault {
		s.repo.ResetDefault(userID)
	}
	addr := &model.Address{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: req.IsDefault,
	}
	if err := s.repo.Create(addr); err != nil {
		return nil, err
	}
	resp := response.ToAddressResponse(addr)
	return &resp, nil
}

func (s *service) List(userID uint) ([]response.AddressResponse, error) {
	addrs, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return response.ToAddressList(addrs), nil
}

func (s *service) Update(userID, addressID uint, req *request.UpdateAddressRequest) (*response.AddressResponse, error) {
	addr, err := s.repo.FindByID(addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	if addr.UserID != userID {
		return nil, ErrAddressNotFound
	}
	if req.IsDefault != nil && *req.IsDefault {
		s.repo.ResetDefault(userID)
	}
	if req.Name != "" {
		addr.Name = req.Name
	}
	if req.Phone != "" {
		addr.Phone = req.Phone
	}
	if req.Province != "" {
		addr.Province = req.Province
	}
	if req.City != "" {
		addr.City = req.City
	}
	if req.District != "" {
		addr.District = req.District
	}
	if req.Detail != "" {
		addr.Detail = req.Detail
	}
	if req.IsDefault != nil {
		addr.IsDefault = *req.IsDefault
	}
	if err := s.repo.Update(addr); err != nil {
		return nil, err
	}
	resp := response.ToAddressResponse(addr)
	return &resp, nil
}

func (s *service) Delete(userID, addressID uint) error {
	addr, err := s.repo.FindByID(addressID)
	if err != nil {
		return ErrAddressNotFound
	}
	if addr.UserID != userID {
		return ErrAddressNotFound
	}
	return s.repo.Delete(addressID, userID)
}
