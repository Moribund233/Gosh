package category

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/category"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create a category
// @Description Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param body body request.CreateCategoryRequest true "Category info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /admin/categories [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create category failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Update a category
// @Description Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body request.UpdateCategoryRequest true "Category update info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /admin/categories/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid category id")
		return
	}
	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update category failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /admin/categories/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid category id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrHasChildren {
			response.ErrorWithCode(c, 409, errcode.ErrConflict, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "delete category failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Get category tree
// @Description Get hierarchical category tree
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /categories [get]
func (h *Handler) Tree(c *gin.Context) {
	tree, err := h.svc.GetTree()
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get category tree failed")
		return
	}
	response.Success(c, tree)
}

// @Summary Get category by ID
// @Description Get category info by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /categories/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid category id")
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		if err == svc.ErrCategoryNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get category failed")
		return
	}
	response.Success(c, resp)
}
