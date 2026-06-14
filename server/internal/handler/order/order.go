package order

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/order"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	idempotentKey := c.GetHeader("Idempotent-Key")
	if idempotentKey == "" {
		response.BadRequest(c, "missing Idempotent-Key header")
		return
	}

	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req, idempotentKey)
	if err != nil {
		if errors.Is(err, svc.ErrCartEmpty) || errors.Is(err, svc.ErrNoDefaultAddress) || errors.Is(err, svc.ErrInsufficientStock) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "create order failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	var req request.ListOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "list orders failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.GetByID(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get order failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	var req request.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Cancel(userID.(uint), uint(id), &req); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "cancel order failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Pay(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "pay order failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Ship(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	if err := h.svc.Ship(uint(id)); err != nil {
		if err == svc.ErrOrderNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "ship order failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Confirm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Confirm(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "confirm order failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) ApplyPoints(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	var req request.ApplyPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.ApplyPoints(userID.(uint), uint(id), &req); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus || err == svc.ErrInsufficientPoints {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "apply points failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Rebuy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Rebuy(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "rebuy failed")
		return
	}
	response.Success(c, nil)
}
