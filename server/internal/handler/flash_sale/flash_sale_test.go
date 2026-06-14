package flash_sale

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/model"
	"gosh/internal/testutil"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}
	cfg := config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"}
	err := database.Init(cfg)
	require.NoError(t, err)
	err = database.DB.AutoMigrate(&model.FlashSale{})
	require.NoError(t, err)
}

func createTestFlashSale(t *testing.T) *model.FlashSale {
	t.Helper()
	fs := &model.FlashSale{
		ProductID:  1,
		SKUID:      1,
		FlashPrice: 1990,
		FlashStock: 100,
		StartAt:    time.Now().Add(-time.Hour),
		EndAt:      time.Now().Add(2 * time.Hour),
		Status:     model.FlashSaleStatusActive,
	}
	database.DB.Create(fs)
	return fs
}

func TestListActive_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	createTestFlashSale(t)

	c, w := testutil.NewGinContext("GET", "/api/v1/flash-sales", "")
	h.ListActive(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, float64(1990), resp.Data[0]["flash_price"])
	assert.True(t, resp.Data[0]["countdown"].(float64) > 0)
}

func TestListActive_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/flash-sales", "")
	h.ListActive(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 0)
}

func TestListActive_Expired(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	fs := &model.FlashSale{
		ProductID:  1,
		SKUID:      1,
		FlashPrice: 1990,
		FlashStock: 100,
		StartAt:    time.Now().Add(-4 * time.Hour),
		EndAt:      time.Now().Add(-2 * time.Hour),
		Status:     model.FlashSaleStatusActive,
	}
	database.DB.Create(fs)

	c, w := testutil.NewGinContext("GET", "/api/v1/flash-sales", "")
	h.ListActive(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 0)
}
