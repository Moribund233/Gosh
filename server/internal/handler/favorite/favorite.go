package favorite

import (
	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/favorite"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Add a favorite
// @Description Add a product to favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.AddFavoriteRequest true "Product ID"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites [post]
func (h *Handler) Add(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Add(userID.(uint), req.ProductID)
	if err != nil {
		if err == svc.ErrFavoriteExists {
			response.ErrorWithCode(c, 409, errcode.ErrConflict, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "add favorite failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Remove a favorite
// @Description Remove a product from favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.AddFavoriteRequest true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/remove [post]
func (h *Handler) Remove(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.Remove(userID.(uint), req.ProductID); err != nil {
		if err == svc.ErrFavoriteNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "remove favorite failed")
		return
	}
	response.Success(c, nil)
}

// @Summary List favorites
// @Description List user's favorite products with pagination
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response "data: {list: []gosh/internal/dto/response.FavoriteResponse, total: int64}"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites [get]
func (h *Handler) List(c *gin.Context) {
	var req request.ListFavoriteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list favorites failed")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
