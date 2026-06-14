package product

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/product"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		response.InternalError(c, "create product failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update product failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	resp, err := h.svc.GetByID(uint(id))
	if err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get product failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) List(c *gin.Context) {
	var req request.ListProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	list, total, err := h.svc.List(&req)
	if err != nil {
		response.InternalError(c, "list products failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "delete product failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	status := c.Param("status")
	if status != "on" && status != "off" {
		response.BadRequest(c, "invalid status, must be on or off")
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	if err := h.svc.UpdateStatus(uint(id), status); err != nil {
		if err == svc.ErrProductNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update status failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) Search(c *gin.Context) {
	var req request.ListProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	list, total, err := h.svc.Search(&req, uid)
	if err != nil {
		response.InternalError(c, "search failed")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *Handler) HotSearch(c *gin.Context) {
	list, err := h.svc.HotSearch()
	if err != nil {
		response.InternalError(c, "get hot search failed")
		return
	}
	response.Success(c, list)
}

func (h *Handler) SearchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	histories, err := h.svc.SearchHistory(uid)
	if err != nil {
		response.InternalError(c, "get search history failed")
		return
	}
	response.Success(c, histories)
}

func (h *Handler) ClearSearchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	if err := h.svc.ClearSearchHistory(uid); err != nil {
		response.InternalError(c, "clear search history failed")
		return
	}
	response.Success(c, nil)
}
