package upload

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gosh/internal/config"
	"gosh/pkg/errcode"
	"gosh/pkg/response"
)

// @Summary Upload a file
// @Description Upload an image file (multipart/form-data)
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /upload [post]
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "file is required")
		return
	}

	maxSize := int64(config.AppConfig.Upload.MaxSize) * 1024 * 1024
	if file.Size > maxSize {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, fmt.Sprintf("file too large, max %dMB", config.AppConfig.Upload.MaxSize))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "unsupported file type, allowed: jpg,jpeg,png,gif,webp")
		return
	}

	uploadDir := config.AppConfig.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create upload dir failed")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "save file failed")
		return
	}

	response.Success(c, gin.H{
		"url":  "/uploads/" + filename,
		"size": file.Size,
	})
}

// @Summary Upload a base64 image
// @Description Upload a base64 encoded image
// @Tags Upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object true "Base64 image data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /upload/base64 [post]
func UploadBase64(c *gin.Context) {
	var req struct {
		Data string `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "base64 data is required")
		return
	}

	parts := strings.SplitN(req.Data, ",", 2)
	data := req.Data
	if len(parts) == 2 {
		data = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		response.BadRequestWithCode(c, errcode.ErrBadRequest, "invalid base64 data")
		return
	}

	uploadDir := config.AppConfig.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "create upload dir failed")
		return
	}

	filename := fmt.Sprintf("%d.png", time.Now().UnixNano())
	dst := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(dst, decoded, 0644); err != nil {
		response.InternalErrorWithCode(c, errcode.ErrInternal, "write file failed")
		return
	}

	response.Success(c, gin.H{
		"url": "/uploads/" + filename,
	})
}
