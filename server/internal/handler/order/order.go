package order

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/order"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create order
// @Description Create a new order from selected cart items
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Idempotent-Key header string true "Idempotency key"
// @Param request body request.CreateOrderRequest true "Order creation details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	idempotentKey := c.GetHeader("Idempotent-Key")
	if idempotentKey == "" {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "missing Idempotent-Key header")
		return
	}

	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req, idempotentKey)
	if err != nil {
		if errors.Is(err, svc.ErrCartEmpty) || errors.Is(err, svc.ErrNoDefaultAddress) || errors.Is(err, svc.ErrInsufficientStock) {
			switch {
			case errors.Is(err, svc.ErrCartEmpty):
				response.BadRequestWithCode(c, errcode.ErrCartEmpty, err.Error())
			case errors.Is(err, svc.ErrNoDefaultAddress):
				response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
			case errors.Is(err, svc.ErrInsufficientStock):
				response.BadRequestWithCode(c, errcode.ErrInsufficientStock, err.Error())
			}
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create order failed")
		return
	}
	response.Created(c, resp)
}

// @Summary List orders
// @Description List current user's orders with pagination and filters
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request query request.ListOrderRequest true "List filters"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [get]
func (h *Handler) List(c *gin.Context) {
	var req request.ListOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list orders failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// @Summary Get order by ID
// @Description Get detailed info of a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.GetByID(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get order failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Cancel order
// @Description Cancel an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body request.CancelOrderRequest true "Cancellation reason"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/cancel [post]
func (h *Handler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	var req request.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Cancel(userID.(uint), uint(id), &req); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "cancel order failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Pay order
// @Description Initiate payment for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/pay [post]
func (h *Handler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Pay(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "pay order failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Ship order (admin)
// @Description Mark an order as shipped (admin only)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/ship [post]
func (h *Handler) Ship(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	if err := h.svc.Ship(uint(id)); err != nil {
		if err == svc.ErrOrderNotFound {
			response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "ship order failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Confirm receipt
// @Description Confirm order receipt (mark as completed)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/confirm [post]
func (h *Handler) Confirm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Confirm(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "confirm order failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Apply points to order
// @Description Apply loyalty points discount to an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body request.ApplyPointsRequest true "Points to apply"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/apply-points [post]
func (h *Handler) ApplyPoints(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	var req request.ApplyPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.ApplyPoints(userID.(uint), uint(id), &req); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		if err == svc.ErrInvalidOrderStatus || err == svc.ErrInsufficientPoints {
			if err == svc.ErrInvalidOrderStatus {
				response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			} else if err == svc.ErrInsufficientPoints {
				response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
			}
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "apply points failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Rebuy order
// @Description Add all items from a previous order back to the cart
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/rebuy [post]
func (h *Handler) Rebuy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Rebuy(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			if err == svc.ErrOrderNotFound {
				response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())
			} else if err == svc.ErrOrderNotBelongToUser {
				response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			}
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequestWithCode(c, errcode.ErrOrderStatus, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "rebuy failed")
		return
	}
	response.Success(c, resp)
}
