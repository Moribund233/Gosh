package address

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/address"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create a new address
// @Description Create a new shipping address
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateAddressRequest true "Address info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /addresses [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create address failed")
		return
	}
	response.Created(c, resp)
}

// @Summary List user addresses
// @Description Get all addresses for current user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /addresses [get]
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	addrs, err := h.svc.List(userID.(uint))
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list addresses failed")
		return
	}
	response.Success(c, addrs)
}

// @Summary Update an address
// @Description Update a shipping address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param body body request.UpdateAddressRequest true "Updated address info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /addresses/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid address id")
		return
	}
	var req request.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Update(userID.(uint), uint(id), &req)
	if err != nil {
		if err == svc.ErrAddressNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update address failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Delete an address
// @Description Delete a shipping address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /addresses/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid address id")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Delete(userID.(uint), uint(id)); err != nil {
		if err == svc.ErrAddressNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "delete address failed")
		return
	}
	response.Success(c, nil)
}
