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
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "ok", resp.Message)
	assert.Contains(t, string(resp.Data), "value")
}

func TestCreated(t *testing.T) {
	c, w := testutil.NewGinContext("POST", "/test", "")
	Created(c, gin.H{"id": 1})
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestErrorWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	ErrorWithCode(c, http.StatusConflict, 4001, "conflict")
	assert.Equal(t, http.StatusConflict, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 4001, resp.Code)
	assert.Equal(t, "conflict", resp.Message)
}

func TestBadRequestWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	BadRequestWithCode(c, 1001, "bad request")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1001, resp.Code)
	assert.Equal(t, "bad request", resp.Message)
}

func TestNotFoundWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	NotFoundWithCode(c, 2001, "not found")
	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 2001, resp.Code)
}

func TestUnauthorizedWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	UnauthorizedWithCode(c, 1002, "unauthorized")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1002, resp.Code)
}

func TestForbiddenWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	ForbiddenWithCode(c, 1003, "forbidden")
	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1003, resp.Code)
}

func TestInternalErrorWithCode(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	InternalErrorWithCode(c, 1005, "internal error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1005, resp.Code)
}

func TestBadRequest(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	BadRequest(c, "bad request")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp TestResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, -1, resp.Code)
}

func TestNotFound(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	NotFound(c, "not found")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUnauthorized(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Unauthorized(c, "unauthorized")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestForbidden(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	Forbidden(c, "forbidden")
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestInternalError(t *testing.T) {
	c, w := testutil.NewGinContext("GET", "/test", "")
	InternalError(c, "internal error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type TestResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}
