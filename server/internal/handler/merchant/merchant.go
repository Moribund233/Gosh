package merchant

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/merchant"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Apply for merchant
// @Description User applies to become a merchant
// @Tags Merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.ApplyMerchantRequest true "Application info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /merchant/apply [post]
func (h *Handler) Apply(c *gin.Context) {
	var req request.ApplyMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Apply(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrApplicationExists {
			response.ErrorWithCode(c, 409, errcode.ErrConflict, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "apply merchant failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Get my application
// @Description Get the current user's merchant application
// @Tags Merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /merchant/application [get]
func (h *Handler) MyApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.svc.MyApplication(userID.(uint))
	if err != nil {
		if err == svc.ErrApplicationNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get application failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Review merchant application
// @Description Admin reviews a merchant application
// @Tags Merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.ReviewMerchantRequest true "Review info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /admin/merchant/review [post]
func (h *Handler) Review(c *gin.Context) {
	var req request.ReviewMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	resp, err := h.svc.Review(0, &req)
	if err != nil {
		if err == svc.ErrApplicationNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		if err == svc.ErrAlreadyReviewed {
			response.ErrorWithCode(c, 409, errcode.ErrConflict, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "review application failed")
		return
	}
	response.Success(c, resp)
}

// @Summary List merchant applications
// @Description Admin lists all merchant applications
// @Tags Merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Application status filter"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/merchant/applications [get]
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
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list applications failed")
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
