package review

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/review"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Create a review
// @Description Create a product review
// @Tags Reviews
// @Accept json
// @Produce json
// @Param body body request.CreateReviewRequest true "Review info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /reviews [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create review failed")
		return
	}
	response.Created(c, resp)
}

// @Summary List reviews
// @Description List reviews for a product with pagination
// @Tags Reviews
// @Accept json
// @Produce json
// @Param product_id query int true "Product ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /reviews [get]
func (h *Handler) List(c *gin.Context) {
	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid product_id")
		return
	}
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page := 1
	size := 10
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 50 {
		size = s
	}
	list, total, err := h.svc.List(uint(productID), page, size)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list reviews failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}
