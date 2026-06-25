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
				fields := []zap.Field{
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				}
				if rid, ok := c.Get(RequestIDKey); ok {
					fields = append(fields, zap.String("request_id", rid.(string)))
				}
				log.Error("panic recovered", fields...)
				response.InternalError(c, "server internal error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
