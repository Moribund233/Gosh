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
	"gosh/pkg/response"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	maxSize := int64(config.AppConfig.Upload.MaxSize) * 1024 * 1024
	if file.Size > maxSize {
		response.BadRequest(c, fmt.Sprintf("file too large, max %dMB", config.AppConfig.Upload.MaxSize))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		response.BadRequest(c, "unsupported file type, allowed: jpg,jpeg,png,gif,webp")
		return
	}

	uploadDir := config.AppConfig.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "create upload dir failed")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.InternalError(c, "save file failed")
		return
	}

	response.Success(c, gin.H{
		"url":  "/uploads/" + filename,
		"size": file.Size,
	})
}

func UploadBase64(c *gin.Context) {
	var req struct {
		Data string `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "base64 data is required")
		return
	}

	parts := strings.SplitN(req.Data, ",", 2)
	data := req.Data
	if len(parts) == 2 {
		data = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		response.BadRequest(c, "invalid base64 data")
		return
	}

	uploadDir := config.AppConfig.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "create upload dir failed")
		return
	}

	filename := fmt.Sprintf("%d.png", time.Now().UnixNano())
	dst := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(dst, decoded, 0644); err != nil {
		response.InternalError(c, "write file failed")
		return
	}

	response.Success(c, gin.H{
		"url": "/uploads/" + filename,
	})
}
