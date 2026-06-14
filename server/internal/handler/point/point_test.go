package point

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
	err = database.DB.AutoMigrate(&model.User{}, &model.PointLog{})
	require.NoError(t, err)
}

func TestGetBalance(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	user := &model.User{Phone: "13800138001", Password: "hash", Nickname: "Test", Points: 500}
	database.DB.Create(user)

	c, w := testutil.NewGinContext("GET", "/api/v1/points", "")
	c.Set("user_id", user.ID)
	h.GetBalance(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp.Data["points"])
}

func TestListLogs_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	user := &model.User{Phone: "13800138001", Password: "hash", Nickname: "Test"}
	database.DB.Create(user)

	c, w := testutil.NewGinContext("GET", "/api/v1/points/logs", "")
	c.Set("user_id", user.ID)
	h.ListLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp.Data["total"])
}

func TestListLogs_WithItems(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	user := &model.User{Phone: "13800138001", Password: "hash", Nickname: "Test", Points: 100}
	database.DB.Create(user)

	database.DB.Create(&model.PointLog{
		UserID:  user.ID,
		Type:    model.PointTypeEarn,
		Amount:  100,
		Balance: 100,
		Note:    "测试积分",
	})

	c, w := testutil.NewGinContext("GET", "/api/v1/points/logs", "")
	c.Set("user_id", user.ID)
	h.ListLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["total"])
}
