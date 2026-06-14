package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gosh/internal/dto/request"
	svc "gosh/internal/service/user"
	"gosh/pkg/response"
)

type Handler struct {
	svc svc.Service
}

func NewHandler() *Handler {
	return &Handler{svc: svc.New()}
}

func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tokenResp, err := h.svc.Register(req.Phone, req.Password, req.Nickname)
	if err != nil {
		if err == svc.ErrPhoneExists {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.InternalError(c, "register failed")
		return
	}
	response.Created(c, tokenResp)
}

func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tokenResp, err := h.svc.Login(req.Phone, req.Password)
	if err != nil {
		if err == svc.ErrInvalidCreds {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, "login failed")
		return
	}
	response.Success(c, tokenResp)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.svc.GetProfile(userID.(uint))
	if err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "get profile failed")
		return
	}
	response.Success(c, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.UpdateProfile(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update profile failed")
		return
	}
	response.Success(c, resp)
}

func (h *Handler) ListUsers(c *gin.Context) {
	var req request.ListRoleUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	users, total, err := h.svc.ListByRole(req.Role, req.Page, req.Size)
	if err != nil {
		response.InternalError(c, "list users failed")
		return
	}
	response.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req request.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateRole(req.UserID, req.Role); err != nil {
		if err == svc.ErrUserNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "update role failed")
		return
	}
	response.Success(c, nil)
}
