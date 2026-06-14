package user

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
	err = database.DB.AutoMigrate(&model.User{})
	require.NoError(t, err)
}

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	c, w := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"测试用户"}`)
	h.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotEmpty(t, resp.Data["token"])
	assert.Equal(t, "13800138000", resp.Data["user"].(map[string]interface{})["phone"])
	assert.Equal(t, "测试用户", resp.Data["user"].(map[string]interface{})["nickname"])
}

func TestRegister_DuplicatePhone(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	body := `{"phone":"13800138001","password":"pass123","nickname":"u1"}`
	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", body)
	h.Register(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/user/register", body)
	h.Register(c2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestRegister_InvalidParams(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	tests := []struct {
		name string
		body string
	}{
		{"missing phone", `{"password":"pass123","nickname":"test"}`},
		{"short phone", `{"phone":"123","password":"pass123","nickname":"test"}`},
		{"missing password", `{"phone":"13800138002","nickname":"test"}`},
		{"short password", `{"phone":"13800138002","password":"123","nickname":"test"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := testutil.NewGinContext("POST", "/api/v1/user/register", tt.body)
			h.Register(c)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138010","password":"pass123","nickname":"login_test"}`)
	h.Register(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/user/login", `{"phone":"13800138010","password":"pass123"}`)
	h.Login(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.NotEmpty(t, resp.Data["token"])
	assert.Equal(t, "login_test", resp.Data["user"].(map[string]interface{})["nickname"])
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138020","password":"pass123","nickname":"u"}`)
	h.Register(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/user/login", `{"phone":"13800138020","password":"wrongpass"}`)
	h.Login(c2)

	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestLogin_NotRegistered(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/user/login", `{"phone":"13800138999","password":"pass123"}`)
	h.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfile(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138030","password":"pass123","nickname":"profile_test"}`)
	h.Register(c1)

	var loginResp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	c2, w2 := testutil.NewGinContext("POST", "/api/v1/user/login", `{"phone":"13800138030","password":"pass123"}`)
	h.Login(c2)
	json.Unmarshal(w2.Body.Bytes(), &loginResp)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/user/profile", "")
	c3.Set("user_id", uint(1))
	h.GetProfile(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138040","password":"pass123","nickname":"old_name"}`)
	h.Register(c1)

	c2, w2 := testutil.NewGinContext("PUT", "/api/v1/user/profile", `{"nickname":"new_name","avatar":"http://example.com/avatar.jpg"}`)
	c2.Set("user_id", uint(1))
	h.UpdateProfile(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "new_name", resp.Data["nickname"])
	assert.Equal(t, "http://example.com/avatar.jpg", resp.Data["avatar"])
}

func TestUpdateProfile_Partial(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/user/register", `{"phone":"13800138041","password":"pass123","nickname":"original"}`)
	h.Register(c1)

	c2, w2 := testutil.NewGinContext("PUT", "/api/v1/user/profile", `{"nickname":"only_name"}`)
	c2.Set("user_id", uint(1))
	h.UpdateProfile(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "only_name", resp.Data["nickname"])
}
