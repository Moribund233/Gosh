package errcode

const (
	ErrBadRequest   = 1001
	ErrUnauthorized = 1002
	ErrForbidden    = 1003
	ErrNotFound     = 1004
	ErrConflict     = 1005
	ErrInternal     = 1999
)

const (
	ErrUserExists    = 2001
	ErrPasswordWrong = 2002
	ErrUserNotFound  = 2003
)

const (
	ErrCategoryNotFound  = 3001
	ErrProductNotFound   = 3002
	ErrSKUNotFound       = 3003
	ErrInsufficientStock = 3004
)

const (
	ErrOrderNotFound = 4001
	ErrOrderStatus   = 4002
	ErrCartEmpty     = 4003
)

const (
	ErrPaymentMethod = 5001
	ErrPaymentFailed = 5002
	ErrRefundFailed  = 5003
)

const (
	ErrCouponNotFound = 6001
	ErrCouponSoldOut  = 6002
	ErrCouponReceived = 6003
)
