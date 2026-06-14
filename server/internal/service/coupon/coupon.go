package coupon

import (
	"errors"
	"time"

	"gosh/internal/dto/request"
	"gosh/internal/dto/response"
	"gosh/internal/model"
	repo "gosh/internal/repository/coupon"
)

var (
	ErrCouponNotFound      = errors.New("coupon not found")
	ErrCouponExpired       = errors.New("coupon has expired")
	ErrCouponNotStarted    = errors.New("coupon not yet started")
	ErrCouponSoldOut       = errors.New("coupon sold out")
	ErrCouponLimitReached  = errors.New("per-user limit reached")
	ErrAlreadyReceived     = errors.New("already received this coupon")
	ErrCouponNotApplicable = errors.New("order amount does not meet coupon condition")
)

type Service interface {
	Create(req *request.CreateCouponRequest) (*response.CouponResponse, error)
	Receive(userID, couponID uint) (*response.UserCouponResponse, error)
	GetAvailable(userID uint, amount int64) ([]response.UserCouponResponse, error)
	Calculate(req *request.CalculateCouponRequest) (*response.CouponCalculateResponse, error)
}

type service struct {
	repo repo.Repository
}

func New() Service {
	return &service{repo: repo.New()}
}

func (s *service) Create(req *request.CreateCouponRequest) (*response.CouponResponse, error) {
	startAt, err := time.Parse("2006-01-02 15:04:05", req.StartAt)
	if err != nil {
		return nil, errors.New("invalid start_at format, use YYYY-MM-DD HH:mm:ss")
	}
	endAt, err := time.Parse("2006-01-02 15:04:05", req.EndAt)
	if err != nil {
		return nil, errors.New("invalid end_at format, use YYYY-MM-DD HH:mm:ss")
	}

	count := req.TotalCount
	if count == 0 {
		count = -1
	}

	coupon := &model.Coupon{
		Name:        req.Name,
		Type:        req.Type,
		Condition:   req.Condition,
		Discount:    req.Discount,
		TotalCount:  req.TotalCount,
		RemainCount: count,
		PerLimit:    req.PerLimit,
		StartAt:     startAt,
		EndAt:       endAt,
		Status:      model.CouponStatusActive,
	}
	if err := s.repo.Create(coupon); err != nil {
		return nil, err
	}
	resp := response.ToCouponResponse(coupon)
	return &resp, nil
}

func (s *service) Receive(userID, couponID uint) (*response.UserCouponResponse, error) {
	coupon, err := s.repo.FindByID(couponID)
	if err != nil {
		return nil, ErrCouponNotFound
	}

	now := time.Now()
	if now.Before(coupon.StartAt) {
		return nil, ErrCouponNotStarted
	}
	if now.After(coupon.EndAt) {
		return nil, ErrCouponExpired
	}
	if coupon.Status != model.CouponStatusActive {
		return nil, ErrCouponExpired
	}

	if coupon.TotalCount > 0 && coupon.RemainCount <= 0 {
		return nil, ErrCouponSoldOut
	}

	existing, err := s.repo.FindUserCoupon(userID, couponID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyReceived
	}

	if err := s.repo.DecrementRemain(couponID); err != nil {
		return nil, ErrCouponSoldOut
	}

	uc := &model.UserCoupon{
		UserID:   userID,
		CouponID: couponID,
		Status:   model.UserCouponStatusUnused,
	}
	if err := s.repo.CreateUserCoupon(uc); err != nil {
		return nil, err
	}

	resp := response.ToUserCouponResponse(uc, coupon)
	return &resp, nil
}

func (s *service) GetAvailable(userID uint, amount int64) ([]response.UserCouponResponse, error) {
	coupons, err := s.repo.FindActive(amount)
	if err != nil {
		return nil, err
	}

	ucs, err := s.repo.ListUserCoupons(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	couponMap := make(map[uint]*model.Coupon)
	for i := range coupons {
		couponMap[coupons[i].ID] = &coupons[i]
	}

	var available []response.UserCouponResponse
	for _, uc := range ucs {
		if uc.Status != model.UserCouponStatusUnused {
			continue
		}
		c, ok := couponMap[uc.CouponID]
		if !ok {
			continue
		}
		if now.Before(c.StartAt) || now.After(c.EndAt) {
			continue
		}

		resp := response.ToUserCouponResponse(&uc, c)
		available = append(available, resp)
	}

	return available, nil
}

func (s *service) Calculate(req *request.CalculateCouponRequest) (*response.CouponCalculateResponse, error) {
	coupon, err := s.repo.FindByID(req.CouponID)
	if err != nil {
		return nil, ErrCouponNotFound
	}

	now := time.Now()
	if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
		return nil, ErrCouponExpired
	}
	if coupon.Status != model.CouponStatusActive {
		return nil, ErrCouponExpired
	}

	if req.OrderAmount < coupon.Condition {
		return nil, ErrCouponNotApplicable
	}

	var discount int64
	switch coupon.Type {
	case model.CouponTypeFullReduce:
		discount = coupon.Discount
		if discount > req.OrderAmount {
			discount = req.OrderAmount
		}
	case model.CouponTypeDiscount:
		discount = req.OrderAmount * (100 - coupon.Discount) / 100
		if discount > req.OrderAmount {
			discount = req.OrderAmount
		}
	}

	payAmount := req.OrderAmount - discount

	return &response.CouponCalculateResponse{
		OriginalAmount: req.OrderAmount,
		DiscountAmount: discount,
		PayAmount:      payAmount,
	}, nil
}
