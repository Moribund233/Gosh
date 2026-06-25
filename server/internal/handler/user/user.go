package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/user"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

// @Summary Register a new user
// @Description Register with phone, password, and optional nickname
// @Tags Users
// @Accept json
// @Produce json
// @Param body body request.RegisterRequest true "Registration info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /user/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	tokenResp, err := h.svc.Register(req.Phone, req.Password, req.Nickname)
	if err != nil {
		if err == svc.ErrPhoneExists {
			response.ErrorWithCode(c, http.StatusConflict, errcode.ErrUserExists, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "register failed")
		return
	}
	response.Created(c, tokenResp)
}

// @Summary User login
// @Description Login with phone and password
// @Tags Users
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	tokenResp, err := h.svc.Login(req.Phone, req.Password)
	if err != nil {
		if err == svc.ErrInvalidCreds {
			response.UnauthorizedWithCode(c, errcode.ErrPasswordWrong, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "login failed")
		return
	}
	response.Success(c, tokenResp)
}

// @Summary Get user profile
// @Description Get current user profile with favorites count and order stats
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.svc.GetProfile(userID.(uint))
	if err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFoundWithCode(c, errcode.ErrUserNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "get profile failed")
		return
	}
	response.Success(c, profile)
}

// @Summary Update user profile
// @Description Update nickname and/or avatar
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.UpdateProfileRequest true "Profile update info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.UpdateProfile(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFoundWithCode(c, errcode.ErrUserNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update profile failed")
		return
	}
	response.Success(c, resp)
}

// @Summary List users by role
// @Description List users filtered by role (super admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role query string false "Role filter"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response "data: {list: []gosh/internal/dto/response.UserResponse, total: int64}"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	var req request.ListRoleUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	users, total, err := h.svc.ListByRole(req.Role, req.Page, req.Size)
	if err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "list users failed")
		return
	}
	response.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

// @Summary Update user role
// @Description Update a user's role (super admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.UpdateUserRoleRequest true "User ID and new role"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/role [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	var req request.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())
		return
	}
	if err := h.svc.UpdateRole(req.UserID, req.Role); err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFoundWithCode(c, errcode.ErrUserNotFound, err.Error())
			return
		}
		response.InternalErrorWithCode(c, errcode.ErrInternal, "update role failed")
		return
	}
	response.Success(c, nil)
}
