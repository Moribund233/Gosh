package favorite

import (
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/favorite"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Add(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Add(userID.(uint), req.ProductID)
	if err != nil {
		if err == svc.ErrFavoriteExists {
			response.Error(c, 409, err.Error())
			return
		}
		response.InternalError(c, "add favorite failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Remove(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Remove(userID.(uint), req.ProductID); err != nil {
		if err == svc.ErrFavoriteNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "remove favorite failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) List(c *gin.Context) {
	var req request.ListFavoriteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalError(c, "list favorites failed")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
