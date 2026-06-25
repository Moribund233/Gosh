package payment

import (
	"io"

	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/payment"
	"gosh/internal/dto/request"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Get payment methods
// @Description Get available payment methods
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payment/methods [get]
func (h *Handler) GetMethods(c *gin.Context) {
	methods := h.svc.GetMethods()
	response.Success(c, methods)
}

// @Summary Initiate payment
// @Description Initiate payment for an order using the specified method
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.PayRequest true "Payment details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payment/pay [post]
func (h *Handler) Pay(c *gin.Context) {
	var req request.PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Pay(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "payment failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Payment callback
// @Description Handle payment gateway callback/notification
// @Tags Payments
// @Accept plain
// @Produce plain
// @Param method path string true "Payment method"
// @Param body body string true "Raw callback body"
// @Success 200 {string} string "ok"
// @Router /payment/callback/{method} [post]
func (h *Handler) Callback(c *gin.Context) {
	method := c.Param("method")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "ok")
		return
	}

	if err := h.svc.ProcessCallback(method, body); err != nil {
		c.String(200, "ok")
		return
	}
	c.String(200, "ok")
}

// @Summary Refund payment (admin)
// @Description Process a refund for a completed payment (admin only)
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.RefundRequest true "Refund details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payment/refund [post]
func (h *Handler) Refund(c *gin.Context) {
	var req request.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	err := h.svc.Refund(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrPaymentNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrRefundExists {
			response.BadRequestWithCode(c, errcode.ErrConflict, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "refund failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Get payment status
// @Description Get payment status by order number
// @Tags Payments
// @Accept json
// @Produce json
// @Param order_no path string true "Order number"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payment/status/{order_no} [get]
func (h *Handler) GetStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	resp, err := h.svc.GetStatus(orderNo)
	if err != nil {
		if err == svc.ErrPaymentNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get payment status failed")
		return
	}
	response.Success(c, resp)
}
