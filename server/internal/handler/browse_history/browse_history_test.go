package browse_history

import (
	"encoding/json"
	"net/http"
	"testing"

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
	err = database.DB.AutoMigrate(&model.BrowseHistory{})
	require.NoError(t, err)
}

func TestAddBrowseHistory_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/browse-history", `{"product_id":100}`)
	c.Set("user_id", uint(1))
	h.Add(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, float64(100), resp.Data["product_id"])
}

func TestListBrowseHistory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/browse-history", `{"product_id":1}`)
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, _ := testutil.NewGinContext("POST", "/api/v1/browse-history", `{"product_id":2}`)
	c2.Set("user_id", uint(1))
	h.Add(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/browse-history", "")
	c3.Set("user_id", uint(1))
	h.List(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w3.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	data := resp.Data["list"].([]interface{})
	assert.Len(t, data, 2)
}

func TestListBrowseHistory_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/browse-history", "")
	c.Set("user_id", uint(99))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	data := resp.Data["list"].([]interface{})
	assert.Len(t, data, 0)
}
