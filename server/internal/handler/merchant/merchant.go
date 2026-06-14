package merchant

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/merchant"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Apply(c *gin.Context) {
	var req request.ApplyMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Apply(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrApplicationExists {
			response.Error(c, 409, err.Error())
			return
		}
		response.InternalError(c, "apply merchant failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) MyApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.svc.MyApplication(userID.(uint))
	if err != nil {
		if err == svc.ErrApplicationNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get application failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Review(c *gin.Context) {
	var req request.ReviewMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Review(0, &req)
	if err != nil {
		if err == svc.ErrApplicationNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrAlreadyReviewed {
			response.Error(c, 409, err.Error())
			return
		}
		response.InternalError(c, "review application failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page := 1
	size := 10
	if p, err := parseInt(pageStr, 1); err == nil {
		page = p
	}
	if s, err := parseInt(sizeStr, 10); err == nil {
		size = s
	}
	list, total, err := h.svc.List(status, page, size)
	if err != nil {
		response.InternalError(c, "list applications failed")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func parseInt(s string, defaultVal int) (int, error) {
	if s == "" {
		return defaultVal, nil
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return defaultVal, err
	}
	return n, nil
}
