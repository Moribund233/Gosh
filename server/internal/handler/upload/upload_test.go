package upload

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gosh/internal/config"
	"gosh/internal/testutil"
)

func setupUploadTest(t *testing.T) {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
		Upload: config.UploadConfig{Dir: "/tmp/test-uploads", MaxSize: 10},
	}
}

func createMultipartContext(t *testing.T, fieldName, fileName string, content []byte) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	assert.NoError(t, err)
	_, err = part.Write(content)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return c, w
}

func TestUpload_Success(t *testing.T) {
	setupUploadTest(t)
	c, w := createMultipartContext(t, "file", "test.png", []byte("fake-image-content"))
	Upload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
		} `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Code)
	assert.Contains(t, resp.Data.URL, "/uploads/")
	assert.True(t, resp.Data.Size > 0)
}

func TestUpload_MissingFile(t *testing.T) {
	setupUploadTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/upload", nil)
	c.Request.Header.Set("Content-Type", "multipart/form-data")
	Upload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "file is required", resp.Message)
}

func TestUpload_UnsupportedType(t *testing.T) {
	setupUploadTest(t)
	c, w := createMultipartContext(t, "file", "malware.exe", []byte("bad"))
	Upload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Message, "unsupported file type")
}

func TestUpload_FileTooLarge(t *testing.T) {
	setupUploadTest(t)
	content := make([]byte, 11*1024*1024)
	c, w := createMultipartContext(t, "file", "big.png", content)
	Upload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Message, "file too large")
}

func TestUploadBase64_Success(t *testing.T) {
	setupUploadTest(t)
	c, w := testutil.NewGinContext("POST", "/upload/base64", `{"data":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="}`)
	UploadBase64(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Code)
	assert.Contains(t, resp.Data.URL, "/uploads/")
}

func TestUploadBase64_MissingData(t *testing.T) {
	setupUploadTest(t)
	c, w := testutil.NewGinContext("POST", "/upload/base64", `{}`)
	UploadBase64(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "base64 data is required", resp.Message)
}

func TestUploadBase64_InvalidData(t *testing.T) {
	setupUploadTest(t)
	c, w := testutil.NewGinContext("POST", "/upload/base64", `{"data":"not-base64!!!"}`)
	UploadBase64(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid base64 data", resp.Message)
}
