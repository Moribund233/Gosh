package coupon

import (
	"time"

	"gorm.io/gorm"
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Create(coupon *model.Coupon) error
	FindByID(id uint) (*model.Coupon, error)
	FindActive(amount int64) ([]model.Coupon, error)
	Update(coupon *model.Coupon) error
	DecrementRemain(id uint) error
	CreateUserCoupon(uc *model.UserCoupon) error
	FindUserCoupon(userID, couponID uint) (*model.UserCoupon, error)
	ListUserCoupons(userID uint) ([]model.UserCoupon, error)
	UpdateUserCoupon(uc *model.UserCoupon) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Create(coupon *model.Coupon) error {
	return database.DB.Create(coupon).Error
}

func (r *repo) FindByID(id uint) (*model.Coupon, error) {
	var c model.Coupon
	err := database.DB.First(&c, id).Error
	return &c, err
}

func (r *repo) FindActive(amount int64) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()
	err := database.DB.Where("status = ? AND start_at <= ? AND end_at > ? AND condition <= ?",
		model.CouponStatusActive, now, now, amount).
		Find(&coupons).Error
	return coupons, err
}

func (r *repo) Update(coupon *model.Coupon) error {
	return database.DB.Save(coupon).Error
}

func (r *repo) DecrementRemain(id uint) error {
	return database.DB.Model(&model.Coupon{}).
		Where("id = ? AND remain_count > 0", id).
		Update("remain_count", gorm.Expr("remain_count - 1")).Error
}

func (r *repo) CreateUserCoupon(uc *model.UserCoupon) error {
	return database.DB.Create(uc).Error
}

func (r *repo) FindUserCoupon(userID, couponID uint) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	err := database.DB.Where("user_id = ? AND coupon_id = ?", userID, couponID).First(&uc).Error
	return &uc, err
}

func (r *repo) ListUserCoupons(userID uint) ([]model.UserCoupon, error) {
	var ucs []model.UserCoupon
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(100).Find(&ucs).Error
	return ucs, err
}

func (r *repo) UpdateUserCoupon(uc *model.UserCoupon) error {
	return database.DB.Save(uc).Error
}
