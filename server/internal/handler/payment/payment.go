package payment

import (
	"io"

	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/payment"
	"gosh/internal/dto/request"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) GetMethods(c *gin.Context) {
	methods := h.svc.GetMethods()
	response.Success(c, methods)
}

func (h *Handler) Pay(c *gin.Context) {
	var req request.PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Pay(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "payment failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Callback(c *gin.Context) {
	method := c.Param("method")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, 200, "ok")
		return
	}

	if err := h.svc.ProcessCallback(method, body); err != nil {
		c.String(200, "ok")
		return
	}
	c.String(200, "ok")
}

func (h *Handler) Refund(c *gin.Context) {
	var req request.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	err := h.svc.Refund(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrPaymentNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrRefundExists {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "refund failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) GetStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	resp, err := h.svc.GetStatus(orderNo)
	if err != nil {
		if err == svc.ErrPaymentNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get payment status failed")
		return
	}
	response.Success(c, resp)
}
