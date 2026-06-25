package point

import (
	"strconv"

	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/point"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Get point balance
// @Description Get the current user's point balance
// @Tags Points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /points [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	balance, err := h.svc.GetBalance(userID.(uint))
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get point balance failed")
		return
	}
	response.Success(c, gin.H{"points": balance})
}

// @Summary Get point logs
// @Description Get the current user's point transaction history
// @Tags Points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /points/logs [get]
func (h *Handler) ListLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	logs, total, err := h.svc.ListLogs(userID.(uint), page, size)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get point logs failed")
		return
	}
	response.Success(c, gin.H{"list": logs, "total": total})
}
