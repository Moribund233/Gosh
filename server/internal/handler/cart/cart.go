package cart

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/cart"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Add(c *gin.Context) {
	var req request.AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Add(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrSKUNotFound || err == svc.ErrProductOffShelf {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "add to cart failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.svc.List(userID.(uint))
	if err != nil {
		response.InternalError(c, "list cart failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart id")
		return
	}
	var req request.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Update(uint(id), userID.(uint), &req)
	if err != nil {
		if err == svc.ErrCartNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update cart failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart id")
		return
	}
	userID, _ := c.Get("user_id")

	if err := h.svc.Delete(uint(id), userID.(uint)); err != nil {
		if err == svc.ErrCartNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "delete cart failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Select(c *gin.Context) {
	var req request.SelectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	if err := h.svc.Select(userID.(uint), &req); err != nil {
		response.InternalError(c, "select failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Merge(c *gin.Context) {
	var req request.MergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Merge(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "merge cart failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Count(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := h.svc.Count(userID.(uint))
	if err != nil {
		response.InternalError(c, "get cart count failed")
		return
	}
	response.Success(c, gin.H{"count": count})
}
