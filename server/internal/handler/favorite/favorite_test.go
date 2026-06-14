package favorite

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
	err = database.DB.AutoMigrate(&model.Favorite{})
	require.NoError(t, err)
}

func TestAddFavorite_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":100}`)
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

func TestAddFavorite_Duplicate(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":100}`)
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":100}`)
	c2.Set("user_id", uint(1))
	h.Add(c2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestRemoveFavorite(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":200}`)
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/favorites/remove", `{"product_id":200}`)
	c2.Set("user_id", uint(1))
	h.Remove(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRemoveFavorite_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/favorites/remove", `{"product_id":999}`)
	c.Set("user_id", uint(1))
	h.Remove(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListFavorites(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":1}`)
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, _ := testutil.NewGinContext("POST", "/api/v1/favorites", `{"product_id":2}`)
	c2.Set("user_id", uint(1))
	h.Add(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/favorites", "")
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
	assert.Equal(t, float64(2), resp.Data["total"])
}
