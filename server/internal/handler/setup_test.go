package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/model"
	"gosh/internal/router"
)

func setupTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
		Upload: config.UploadConfig{Dir: "/tmp/test-uploads", MaxSize: 10},
	}
	database.Init(config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"})
	database.DB.AutoMigrate(
		&model.User{},
		&model.Address{},
		&model.Favorite{},
		&model.BrowseHistory{},
		&model.MerchantApplication{},
		&model.Category{},
		&model.Product{},
		&model.ProductSKU{},
		&model.ProductReview{},
		&model.SearchHistory{},
		&model.HotSearch{},
		&model.Banner{},
		&model.BrandStory{},
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderLog{},
		&model.IdempotencyRecord{},
		&model.Payment{},
		&model.Coupon{},
		&model.UserCoupon{},
		&model.FlashSale{},
		&model.PointLog{},
	)
	log, _ := zap.NewDevelopment()
	return router.New(log)
}
