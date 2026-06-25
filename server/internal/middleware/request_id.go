package middleware

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x%x", time.Now().UnixMicro(), b)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = generateRequestID()
		}
		c.Set(RequestIDKey, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}
