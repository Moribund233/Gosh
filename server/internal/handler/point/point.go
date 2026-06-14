package point

import (
	"strconv"

	"github.com/gin-gonic/gin"
	svc "gosh/internal/service/point"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) GetBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	balance, err := h.svc.GetBalance(userID.(uint))
	if err != nil {
		response.InternalError(c, "get point balance failed")
		return
	}
	response.Success(c, gin.H{"points": balance})
}

func (h *Handler) ListLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	logs, total, err := h.svc.ListLogs(userID.(uint), page, size)
	if err != nil {
		response.InternalError(c, "get point logs failed")
		return
	}
	response.Success(c, gin.H{"list": logs, "total": total})
}
