package address

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/address"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "create address failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	addrs, err := h.svc.List(userID.(uint))
	if err != nil {
		response.InternalError(c, "list addresses failed")
		return
	}
	response.Success(c, addrs)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid address id")
		return
	}
	var req request.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Update(userID.(uint), uint(id), &req)
	if err != nil {
		if err == svc.ErrAddressNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update address failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid address id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Delete(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrAddressNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "delete address failed")
		return
	}
	response.Success(c, nil)
}
