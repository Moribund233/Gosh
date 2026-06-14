package testutil

import (
	"encoding/json"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewGinContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func NewGinContextWithToken(method, path, body, token string) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := NewGinContext(method, path, body)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	return c, w
}

func PerformRequest(r *gin.Engine, method, path, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}
	r.ServeHTTP(w, req)
	return w
}

type TestResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}
