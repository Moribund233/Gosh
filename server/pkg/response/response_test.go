package response

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gosh/internal/testutil"
)

func TestSuccess(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Success(c, gin.H{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    *json.RawMessage `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "ok", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestCreated(t *testing.T) {
	c, w := testutil.NewGinContext("POST", "/test", "")
	Created(c, gin.H{"id": 1})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestError(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Error(c, http.StatusBadRequest, "bad request")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBadRequest(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	BadRequest(c, "invalid params")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnauthorized(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Unauthorized(c, "no token")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestForbidden(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Forbidden(c, "no permission")
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestNotFound(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	NotFound(c, "resource not found")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInternalError(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	InternalError(c, "server error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
