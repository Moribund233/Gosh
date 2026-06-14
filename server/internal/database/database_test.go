package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gosh/internal/config"
)

func TestInit_SQLite(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"}
	err := Init(cfg)
	require.NoError(t, err)
	require.NotNil(t, DB)

	sqlDB, err := DB.DB()
	require.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
	sqlDB.Close()
}

func TestInit_Unsupported(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "mysql", Path: "test"}
	err := Init(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database driver")
}

func TestSQLMock(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT current_database()`).WillReturnRows(sqlmock.NewRows([]string{"current_database()"}).AddRow("gosh"))
	var name string
	db.Raw("SELECT current_database()").Scan(&name)
	assert.Equal(t, "gosh", name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
