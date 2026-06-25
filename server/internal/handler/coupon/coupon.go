package coupon

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/coupon"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create a coupon
// @Description Admin creates a new coupon
// @Tags Coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateCouponRequest true "Coupon info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /coupons [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		if err.Error() == "invalid start_at format, use YYYY-MM-DD HH:mm:ss" || err.Error() == "invalid end_at format, use YYYY-MM-DD HH:mm:ss" {
			response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create coupon failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Receive a coupon
// @Description User receives a coupon by ID
// @Tags Coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Coupon ID"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /coupons/{id}/receive [post]
func (h *Handler) Receive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid coupon id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Receive(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrCouponNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrCouponExpired || err == svc.ErrCouponNotStarted || err == svc.ErrCouponSoldOut || err == svc.ErrCouponLimitReached || err == svc.ErrAlreadyReceived {
			response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "receive coupon failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Get available coupons
// @Description Get available coupons for the current user
// @Tags Coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param amount query int false "Order amount for filtering"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /coupons/available [get]
func (h *Handler) GetAvailable(c *gin.Context) {
	userID, _ := c.Get("user_id")
	amountStr := c.Query("amount")
	amount, _ := strconv.ParseInt(amountStr, 10, 64)

	list, err := h.svc.GetAvailable(userID.(uint), amount)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get available coupons failed")
		return
	}
	response.Success(c, list)
}

// @Summary Calculate coupon discount
// @Description Calculate the discount for a coupon
// @Tags Coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CalculateCouponRequest true "Calculation info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /coupons/calculate [post]
func (h *Handler) Calculate(c *gin.Context) {
	var req request.CalculateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Calculate(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrCouponNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrCouponExpired || err == svc.ErrCouponNotApplicable {
			response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "calculate failed")
		return
	}
	response.Success(c, resp)
}
