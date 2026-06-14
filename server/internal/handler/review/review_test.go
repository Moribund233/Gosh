package review

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
	err = database.DB.AutoMigrate(&model.ProductReview{}, &model.User{})
	require.NoError(t, err)
}

func TestCreateReview_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	database.DB.Create(&model.User{Phone: "13800138000", Password: "hash", Nickname: "测试用户", Role: model.RoleUser, Status: model.StatusActive})

	body := `{"product_id":1,"score":5,"content":"很好吃，回购多次了"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/reviews", body)
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, float64(5), resp.Data["score"])
	assert.Equal(t, "很好吃，回购多次了", resp.Data["content"])
	assert.Equal(t, "测***", resp.Data["nickname"])
}

func TestCreateReview_InvalidScore(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"product_id":1,"score":0,"content":"差评"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/reviews", body)
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListReviews(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	database.DB.Create(&model.User{Phone: "13800138001", Password: "hash", Nickname: "用户A", Role: model.RoleUser, Status: model.StatusActive})

	c1, _ := testutil.NewGinContext("POST", "/api/v1/reviews", `{"product_id":1,"score":5,"content":"好评"}`)
	c1.Set("user_id", uint(1))
	h.Create(c1)

	c2, _ := testutil.NewGinContext("POST", "/api/v1/reviews", `{"product_id":1,"score":4,"content":"还不错"}`)
	c2.Set("user_id", uint(1))
	h.Create(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/reviews?product_id=1", "")
	h.List(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w3.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, float64(2), resp.Data["total"])
}

func TestListReviews_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/reviews?product_id=999", "")
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp.Data["total"])
}
