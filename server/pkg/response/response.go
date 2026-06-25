package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    -1,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func ErrorWithCode(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequestWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusBadRequest, code, message)
}

func NotFoundWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusNotFound, code, message)
}

func UnauthorizedWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusUnauthorized, code, message)
}

func ForbiddenWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusForbidden, code, message)
}

func InternalErrorWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusInternalServerError, code, message)
}
