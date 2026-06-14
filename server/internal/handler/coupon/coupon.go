package coupon

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/coupon"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		if err.Error() == "invalid start_at format, use YYYY-MM-DD HH:mm:ss" || err.Error() == "invalid end_at format, use YYYY-MM-DD HH:mm:ss" {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "create coupon failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Receive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid coupon id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Receive(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrCouponNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrCouponExpired || err == svc.ErrCouponNotStarted || err == svc.ErrCouponSoldOut || err == svc.ErrCouponLimitReached || err == svc.ErrAlreadyReceived {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "receive coupon failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) GetAvailable(c *gin.Context) {
	userID, _ := c.Get("user_id")
	amountStr := c.Query("amount")
	amount, _ := strconv.ParseInt(amountStr, 10, 64)

	list, err := h.svc.GetAvailable(userID.(uint), amount)
	if err != nil {
		response.InternalError(c, "get available coupons failed")
		return
	}
	response.Success(c, list)
}

func (h *Handler) Calculate(c *gin.Context) {
	var req request.CalculateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Calculate(&req)
	if err != nil {
		if err == svc.ErrCouponNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrCouponExpired || err == svc.ErrCouponNotApplicable {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "calculate failed")
		return
	}
	response.Success(c, resp)
}
