package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/model"
	"gosh/internal/router"
	"gosh/internal/scheduler"
)

func main() {
	cfgPath := "config/config.yaml"
	if err := config.Init(cfgPath); err != nil {
		log.Fatalf("init config: %v", err)
	}

	logger := initLogger()

	if err := database.Init(config.AppConfig.Database); err != nil {
		logger.Fatal("init database failed", zap.Error(err))
	}

	if err := database.DB.AutoMigrate(
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
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderLog{},
		&model.IdempotencyRecord{},
		&model.Banner{},
		&model.BrandStory{},
		&model.Payment{},
		&model.Coupon{},
		&model.UserCoupon{},
		&model.FlashSale{},
		&model.PointLog{},
	); err != nil {
		logger.Fatal("auto migrate failed", zap.Error(err))
	}
	logger.Info("database migration completed")

	orderScheduler := scheduler.New()
	orderScheduler.Start(logger)

	r := router.New(logger)
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	logger.Info("server starting", zap.String("addr", addr))

	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatal("server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("shutting down server", zap.String("signal", sig.String()))
	orderScheduler.Stop()
	logger.Info("server exited")
}

func initLogger() *zap.Logger {
	lc := config.AppConfig.Logger

	var level zapcore.Level
	if err := level.Set(lc.Level); err != nil {
		level = zapcore.InfoLevel
	}

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	if config.AppConfig.Server.Mode == "debug" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	}

	if lc.Filename != "" {
		if err := os.MkdirAll("logs", 0755); err == nil {
			file, err := os.OpenFile(lc.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(file), level))
			}
		}
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller())
}
