package product

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/product"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create a product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param body body request.CreateProductRequest true "Product info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create product failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Update a product
// @Description Update an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param body body request.UpdateProductRequest true "Product update info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid product id")
		return
	}
	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFoundWithCode(c, errcode.ErrProductNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update product failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Get product by ID
// @Description Get detailed product info including SKUs
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /products/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid product id")
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFoundWithCode(c, errcode.ErrProductNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get product failed")
		return
	}
	response.Success(c, resp)
}

// @Summary List products
// @Description List products with filters and pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param category_id query int false "Category ID"
// @Param tag query string false "Tag filter"
// @Param keyword query string false "Search keyword"
// @Param sort query string false "Sort field (sales, price_newest, price_oldest, newest)"
// @Param status query string false "Status (on, off)"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [get]
func (h *Handler) List(c *gin.Context) {
	var req request.ListProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	list, total, err := h.svc.List(&req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list products failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid product id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFoundWithCode(c, errcode.ErrProductNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "delete product failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Update product status
// @Description Update product status to on or off
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param status path string true "Status (on, off)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products/{id}/status/{status} [put]
func (h *Handler) UpdateStatus(c *gin.Context) {
	status := c.Param("status")
	if status != "on" && status != "off" {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid status, must be on or off")
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid product id")
		return
	}
	if err := h.svc.UpdateStatus(uint(id), status); err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFoundWithCode(c, errcode.ErrProductNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update status failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Search products
// @Description Search products with keyword and filters
// @Tags Products
// @Accept json
// @Produce json
// @Param category_id query int false "Category ID"
// @Param tag query string false "Tag filter"
// @Param keyword query string false "Search keyword"
// @Param sort query string false "Sort field (sales, price_newest, price_oldest, newest)"
// @Param status query string false "Status (on, off)"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/search [get]
func (h *Handler) Search(c *gin.Context) {
	var req request.ListProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	list, total, err := h.svc.Search(&req, uid)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "search failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// @Summary Get hot search keywords
// @Description Get popular search keywords
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/hot-search [get]
func (h *Handler) HotSearch(c *gin.Context) {
	list, err := h.svc.HotSearch()
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get hot search failed")
		return
	}
	response.Success(c, list)
}

// @Summary Get search history
// @Description Get current user's search history
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products/search-history [get]
func (h *Handler) SearchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	histories, err := h.svc.SearchHistory(uid)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get search history failed")
		return
	}
	response.Success(c, histories)
}

// @Summary Clear search history
// @Description Clear current user's search history
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /products/search-history/clear [post]
func (h *Handler) ClearSearchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	if err := h.svc.ClearSearchHistory(uid); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "clear search history failed")
		return
	}
	response.Success(c, nil)
}
