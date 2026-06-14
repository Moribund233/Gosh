package banner

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
		Upload: config.UploadConfig{Dir: "/tmp/test-uploads", MaxSize: 10},
	}
	cfg := config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"}
	err := database.Init(cfg)
	require.NoError(t, err)
	err = database.DB.AutoMigrate(&model.Banner{}, &model.BrandStory{})
	require.NoError(t, err)
}

func TestCreateBanner_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"title":"夏季时令上新","subtitle":"WARM FOOD COLLECTION","description":"产地直送 满99减20","sort_order":1}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/banners", body)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "夏季时令上新", resp.Data["title"])
	assert.Equal(t, "on", resp.Data["status"])
}

func TestListBanners_Active(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"横幅1","sort_order":1}`)
	h.Create(c1)

	c2, _ := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"横幅2","sort_order":2}`)
	h.Create(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/banners", "")
	h.GetActive(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
	var resp struct {
		Code int                      `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w3.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 2)
}

func TestUpdateBanner(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"旧标题","sort_order":1}`)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	body := fmt.Sprintf(`{"title":"新标题","status":"off"}`)
	c2, w2 := testutil.NewGinContext("PUT", fmt.Sprintf("/api/v1/admin/banners/%.0f", createResp.Data.ID), body)
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.Update(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "新标题", resp.Data["title"])
	assert.Equal(t, "off", resp.Data["status"])
}

func TestDeleteBanner(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"待删除","sort_order":1}`)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	c2, w2 := testutil.NewGinContext("DELETE", fmt.Sprintf("/api/v1/admin/banners/%.0f", createResp.Data.ID), "")
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.Delete(c2)

	assert.Equal(t, http.StatusOK, w2.Code)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/admin/banners", "")
	h.List(c3)
	var listResp struct {
		Code int          `json:"code"`
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w3.Body.Bytes(), &listResp)
	assert.Len(t, listResp.Data, 0)
}

func TestGetActive_OnlyOn(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"启用的","sort_order":1}`)
	h.Create(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/banners", `{"title":"禁用的","sort_order":2}`)
	h.Create(c2)
	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &createResp)

	c3, _ := testutil.NewGinContext("PUT", fmt.Sprintf("/api/v1/admin/banners/%.0f", createResp.Data.ID), `{"status":"off"}`)
	c3.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.Update(c3)

	c4, w4 := testutil.NewGinContext("GET", "/api/v1/banners", "")
	h.GetActive(c4)
	var resp struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w4.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 1)
}

func TestBannerEmptyList(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/banners", "")
	h.GetActive(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 0)
}
