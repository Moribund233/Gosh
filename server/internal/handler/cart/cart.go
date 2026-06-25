package cart

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/cart"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Add item to cart
// @Description Add a product SKU to the current user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddCartRequest true "Cart item details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart [post]
func (h *Handler) Add(c *gin.Context) {
	var req request.AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Add(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrSKUNotFound || err == svc.ErrProductOffShelf {
			if err == svc.ErrSKUNotFound {
				response.BadRequestWithCode(c, errcode.ErrSKUNotFound, err.Error())
			} else if err == svc.ErrProductOffShelf {
				response.BadRequestWithCode(c, errcode.ErrProductNotFound, err.Error())
			}
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "add to cart failed")
		return
	}
	response.Created(c, resp)
}

// @Summary Get cart list
// @Description Get current user's cart items with SKU details
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart [get]
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.svc.List(userID.(uint))
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list cart failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Update cart item
// @Description Update quantity or properties of a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart item ID"
// @Param request body request.UpdateCartRequest true "Update details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid cart id")
		return
	}
	var req request.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Update(uint(id), userID.(uint), &req)
	if err != nil {
		if err == svc.ErrCartNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update cart failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Delete cart item
// @Description Remove a product from the current user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart item ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid cart id")
		return
	}
	userID, _ := c.Get("user_id")

	if err := h.svc.Delete(uint(id), userID.(uint)); err != nil {
		if err == svc.ErrCartNotFound {
			response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "delete cart failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Select/deselect cart items
// @Description Select or deselect items in the cart for checkout
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SelectRequest true "Selection details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart/select [post]
func (h *Handler) Select(c *gin.Context) {
	var req request.SelectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	if err := h.svc.Select(userID.(uint), &req); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "select failed")
		return
	}
	response.Success(c, nil)
}

// @Summary Merge guest cart
// @Description Merge guest cart items into the current user's cart after login
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.MergeCartRequest true "Guest cart data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart/merge [post]
func (h *Handler) Merge(c *gin.Context) {
	var req request.MergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")

	resp, err := h.svc.Merge(userID.(uint), &req)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "merge cart failed")
		return
	}
	response.Success(c, resp)
}

// @Summary Get cart item count
// @Description Get the total number of items in the current user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cart/count [get]
func (h *Handler) Count(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := h.svc.Count(userID.(uint))
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get cart count failed")
		return
	}
	response.Success(c, gin.H{"count": count})
}
