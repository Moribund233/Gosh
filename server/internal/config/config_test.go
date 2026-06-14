package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	err := Init("testdata/config.yaml")
	require.NoError(t, err)
	require.NotNil(t, AppConfig)

	assert.Equal(t, 0, AppConfig.Server.Port)
	assert.Equal(t, "test", AppConfig.Server.Mode)
	assert.Equal(t, "sqlite", AppConfig.Database.Driver)
	assert.Equal(t, ":memory:", AppConfig.Database.Path)
	assert.Equal(t, "test-secret", AppConfig.JWT.Secret)
	assert.Equal(t, 72, AppConfig.JWT.ExpireHour)
}

func TestInit_FileNotFound(t *testing.T) {
	err := Init("testdata/nonexistent.yaml")
	assert.Error(t, err)
}

func TestDatabaseDSN_SQLite(t *testing.T) {
	cfg := DatabaseConfig{Driver: "sqlite", Path: "test.db"}
	assert.Equal(t, "test.db", cfg.DSN())
}

func TestDatabaseDSN_Postgres(t *testing.T) {
	cfg := DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "pass",
		DBName:   "gosh",
		SSLMode:  "disable",
	}
	expected := "host=localhost port=5432 user=admin password=pass dbname=gosh sslmode=disable TimeZone=Asia/Shanghai"
	assert.Equal(t, expected, cfg.DSN())
}
