package category

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/category"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "create category failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}
	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update category failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrHasChildren {
			response.Error(c, 409, err.Error())
			return
		}
		response.InternalError(c, "delete category failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Tree(c *gin.Context) {
	tree, err := h.svc.GetTree()
	if err != nil {
		response.InternalError(c, "get category tree failed")
		return
	}
	response.Success(c, tree)
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get category failed")
		return
	}
	response.Success(c, resp)
}
