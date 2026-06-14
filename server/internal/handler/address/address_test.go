package address

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
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
	err = database.DB.AutoMigrate(&model.Address{})
	require.NoError(t, err)
}

func createTestAddress(t *testing.T, h *Handler, userID uint) uint {
	t.Helper()
	body := `{"name":"张三","phone":"13800138000","province":"浙江省","city":"杭州市","district":"西湖区","detail":"文三路138号"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/addresses", body)
	c.Set("user_id", userID)
	h.Create(c)

	var resp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return uint(resp.Data.ID)
}

func TestCreateAddress_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"name":"张三","phone":"13800138000","province":"浙江省","city":"杭州市","district":"西湖区","detail":"文三路138号","is_default":true}`
	c, w := testutil.NewGinContext("POST", "/api/v1/addresses", body)
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "张三", resp.Data["name"])
	assert.Equal(t, "浙江省", resp.Data["province"])
	assert.Equal(t, true, resp.Data["is_default"])
}

func TestCreateAddress_InvalidParams(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/addresses", `{"name":"","phone":"123"}`)
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListAddresses(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	createTestAddress(t, h, 1)
	createTestAddress(t, h, 1)

	c, w := testutil.NewGinContext("GET", "/api/v1/addresses", "")
	c.Set("user_id", uint(1))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                      `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Len(t, resp.Data, 2)
}

func TestListAddresses_OtherUser(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	createTestAddress(t, h, 1)

	c, w := testutil.NewGinContext("GET", "/api/v1/addresses", "")
	c.Set("user_id", uint(2))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                      `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Len(t, resp.Data, 0)
}

func TestUpdateAddress(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	addrID := createTestAddress(t, h, 1)

	c, w := testutil.NewGinContext("PUT", "/api/v1/addresses/1", `{"name":"张四","detail":"新地址"}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", addrID)}}
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
}

func TestDeleteAddress(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	addrID := createTestAddress(t, h, 1)

	c, w := testutil.NewGinContext("DELETE", "/api/v1/addresses/1", "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", addrID)}}
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)

	c2, w2 := testutil.NewGinContext("GET", "/api/v1/addresses", "")
	c2.Set("user_id", uint(1))
	h.List(c2)
	var resp struct {
		Code int                      `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 0)
}

func TestDeleteAddress_OtherUser(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	addrID := createTestAddress(t, h, 1)

	c, w := testutil.NewGinContext("DELETE", "/api/v1/addresses/1", "")
	c.Set("user_id", uint(2))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", addrID)}}
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
