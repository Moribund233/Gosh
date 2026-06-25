package banner

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	"gosh/internal/model"
	svc "gosh/internal/service/banner"
	brandSvc "gosh/internal/service/brand_story"
	"gosh/pkg/errcode"
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

// @Summary Create a banner
// @Description Admin creates a new banner
// @Tags Banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateBannerRequest true "Banner info"
// @Success 201 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/banners [post]
func (h *Handler) Create(c *gin.Context) {
	var req request.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Create(&req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create banner failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Update a banner
// @Description Admin updates a banner
// @Tags Banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Banner ID"
// @Param body body request.UpdateBannerRequest true "Banner update info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/banners/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid banner id")
		return
	}
	var req request.UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Update(uint(id), &req)
	if err != nil {
		if err == svc.ErrBannerNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update banner failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Delete a banner
// @Description Admin deletes a banner
// @Tags Banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Banner ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/banners/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid banner id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		if err == svc.ErrBannerNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "delete banner failed")
		return
	}
	response.Success(c, nil)
}

// @Summary List all banners
// @Description Admin lists all banners
// @Tags Banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Status filter"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/banners [get]
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	list, err := h.svc.List(status)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list banners failed")
		return
	}
	response.Success(c, list)
}

// @Summary Get active banners
// @Description Get all active banners for public display
// @Tags Banners
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /banners [get]
func (h *Handler) GetActive(c *gin.Context) {
	list, err := h.svc.List(model.StatusOn)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get banners failed")
		return
	}
	response.Success(c, list)
}

// @Summary Get brand story
// @Description Get the brand story content
// @Tags Banners
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /brand-story [get]
func (h *Handler) GetBrandStory(c *gin.Context) {
	story, err := h.brandSvc.Get()
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, story)
}

// @Summary Update brand story
// @Description Admin updates the brand story
// @Tags Banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.UpdateBrandStoryRequest true "Brand story info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/brand-story [put]
func (h *Handler) UpdateBrandStory(c *gin.Context) {
	var req request.UpdateBrandStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	story, err := h.brandSvc.Update(&req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update brand story failed")
		return
	}
	response.Success(c, story)
}
