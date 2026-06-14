package banner

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	"gosh/internal/model"
	svc "gosh/internal/service/banner"
	brandSvc "gosh/internal/service/brand_story"
	"gosh/pkg/response"
)

type Handler struct {
	svc      svc.Service
	brandSvc brandSvc.Service
}

func NewHandler() *Handler {
	return &Handler{
		svc:      svc.New(),
		brandSvc: brandSvc.New(),
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		response.InternalError(c, "create banner failed")
		return
	}
	response.Created(c, resp)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}
	var req request.UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrBannerNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update banner failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrBannerNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "delete banner failed")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	list, err := h.svc.List(status)
	if err != nil {
		response.InternalError(c, "list banners failed")
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetActive(c *gin.Context) {
	list, err := h.svc.List(model.StatusOn)
	if err != nil {
		response.InternalError(c, "get banners failed")
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetBrandStory(c *gin.Context) {
	story, err := h.brandSvc.Get()
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, story)
}

func (h *Handler) UpdateBrandStory(c *gin.Context) {
	var req request.UpdateBrandStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	story, err := h.brandSvc.Update(&req)
	if err != nil {
		response.InternalError(c, "update brand story failed")
		return
	}
	response.Success(c, story)
}
