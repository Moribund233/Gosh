package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gosh/internal/config"
)

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

	lvl := logger.Silent
	if config.AppConfig != nil && config.AppConfig.Server.Mode == "debug" {
		lvl = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(lvl),
	})
	if err != nil {
		return fmt.Errorf("connect database failed: %w", err)
	}

	DB = db
	log.Printf("database connected: %s", cfg.Driver)
	return nil
}
