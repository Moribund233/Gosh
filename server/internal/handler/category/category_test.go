package category

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	err = database.DB.AutoMigrate(&model.Category{})
	require.NoError(t, err)
}

func TestCreateCategory_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"name":"方便速食","icon":"🍜","sort_order":1}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/categories", body)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "方便速食", resp.Data["name"])
	assert.Equal(t, "🍜", resp.Data["icon"])
	assert.Equal(t, float64(0), resp.Data["level"])
}

func TestCreateCategory_WithParent(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/categories", `{"name":"食品","sort_order":1}`)
	h.Create(c1)

	var parentResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &parentResp)

	body := fmt.Sprintf(`{"name":"方便速食","parent_id":%.0f,"sort_order":1}`, parentResp.Data.ID)
	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/categories", body)
	h.Create(c2)

	assert.Equal(t, http.StatusCreated, w2.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["level"])
}

func TestCreateCategory_ParentNotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"name":"子分类","parent_id":999,"sort_order":1}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/categories", body)
	h.Create(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetCategoryTree(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, _ := createCtx(`{"name":"食品","sort_order":1}`); h.Create(c1)
	c2, _ := createCtx(`{"name":"生鲜","sort_order":2}`); h.Create(c2)
	c3, _ := createCtx(`{"name":"方便速食","parent_id":1,"sort_order":1}`); h.Create(c3)
	c4, _ := createCtx(`{"name":"新鲜水果","parent_id":2,"sort_order":1}`); h.Create(c4)

	c, w := testutil.NewGinContext("GET", "/api/v1/categories", "")
	h.Tree(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                          `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Len(t, resp.Data, 2)
	children := resp.Data[0]["children"].([]interface{})
	assert.Len(t, children, 1)
}

func TestUpdateCategory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/categories", `{"name":"食品","sort_order":1}`)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	body := fmt.Sprintf(`{"name":"更新名称"}`)
	c2, w2 := testutil.NewGinContext("PUT", "/api/v1/admin/categories/1", body)
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.Update(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestDeleteCategory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/categories", `{"name":"食品","sort_order":1}`)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	c2, w2 := testutil.NewGinContext("DELETE", fmt.Sprintf("/api/v1/admin/categories/%.0f", createResp.Data.ID), "")
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.Delete(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestDeleteCategory_HasChildren(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/categories", `{"name":"食品","sort_order":1}`)
	h.Create(c1)

	var parentResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &parentResp)

	cSub, _ := createCtx(fmt.Sprintf(`{"name":"子分类","parent_id":%.0f,"sort_order":1}`, parentResp.Data.ID)); h.Create(cSub)

	c2, w2 := testutil.NewGinContext("DELETE", fmt.Sprintf("/api/v1/admin/categories/%.0f", parentResp.Data.ID), "")
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", parentResp.Data.ID)}}
	h.Delete(c2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestEmptyCategoryList(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/categories", "")
	h.Tree(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                          `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 0)
}

func createCtx(body string) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/categories", body)
	return c, w
}
