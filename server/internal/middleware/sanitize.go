package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gosh/pkg/sanitize"
)

func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPut &&
			c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body.Close()

		if len(body) == 0 {
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Next()
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Next()
			return
		}

		sanitize.MapStrings(data)

		sanitized, err := json.Marshal(data)
		if err != nil {
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(sanitized))
		c.Request.ContentLength = int64(len(sanitized))
		c.Next()
	}
}
