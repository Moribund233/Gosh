package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Duration("cost", cost),
			zap.String("ip", c.ClientIP()),
		}

		if rid, ok := c.Get(RequestIDKey); ok {
			fields = append(fields, zap.String("request_id", rid.(string)))
		}
		if uid, ok := c.Get("user_id"); ok {
			fields = append(fields, zap.Uint("user_id", uid.(uint)))
		}

		switch {
		case status >= 500:
			log.Error("server error", fields...)
		case status >= 400:
			log.Warn("client error", fields...)
		default:
			log.Info("request", fields...)
		}
	}
}
