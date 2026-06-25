package browse_history

import (
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/browse_history"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Add browse history
// @Description Record a product view in browse history
// @Tags BrowseHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.AddBrowseHistoryRequest true "Product ID"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /browse-history [post]
func (h *Handler) Add(c *gin.Context) {
	var req request.AddBrowseHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Add(userID.(uint), req.ProductID)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "add browse history failed")
		return
	}
	response.Created(c, resp)
}

// @Summary List browse history
// @Description List user's browse history with pagination
// @Tags BrowseHistory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response "data: {list: []gosh/internal/dto/response.BrowseHistoryResponse, total: int64}"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /browse-history [get]
func (h *Handler) List(c *gin.Context) {
	var req request.ListBrowseHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list browse history failed")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
