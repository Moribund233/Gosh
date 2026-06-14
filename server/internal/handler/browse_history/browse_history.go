package browse_history

import (
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/browse_history"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Add(c *gin.Context) {
	var req request.AddBrowseHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Add(userID.(uint), req.ProductID)
	if err != nil {
		response.InternalError(c, "add browse history failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	var req request.ListBrowseHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "list browse history failed")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
