package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gosh/internal/config"
)

const SlowThreshold = 200 * time.Millisecond

const DefaultTimeout = 5 * time.Second

func WithTimeout(d time.Duration) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return DB.WithContext(ctx), cancel
}

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) error {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN())
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	isDebug := config.AppConfig != nil && config.AppConfig.Server.Mode == "debug"

	lvl := logger.Warn
	if isDebug {
		lvl = logger.Info
	}

	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             SlowThreshold,
			LogLevel:                  lvl,
			IgnoreRecordNotFoundError: true,
			Colorful:                  isDebug,
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return fmt.Errorf("connect database failed: %w", err)
	}

	DB = db
	log.Printf("database connected: %s", cfg.Driver)
	return nil
}
