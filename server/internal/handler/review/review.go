package review

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/review"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Create(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "create review failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product_id")
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
		response.InternalError(c, "list reviews failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}
