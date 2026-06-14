package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gosh/pkg/response"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				response.InternalError(c, "server internal error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
