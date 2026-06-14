package product

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
	err = database.DB.AutoMigrate(&model.Product{}, &model.ProductSKU{}, &model.SearchHistory{}, &model.HotSearch{}, &model.Category{})
	require.NoError(t, err)
}

func createTestCategory(t *testing.T) uint {
	t.Helper()
	cat := &model.Category{Name: "测试分类"}
	database.DB.Create(cat)
	return cat.ID
}

func TestCreateProduct_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"有机五常大米","subtitle":"稻花香2号","price":4990,"original_price":6990,"tags":"热卖,有机","skus":[{"name":"5kg/袋","price":4990,"stock":100}]}`, catID)
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "有机五常大米", resp.Data["name"])
	assert.Equal(t, float64(4990), resp.Data["price"])
}

func TestCreateProduct_MissingName(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"category_id":1,"price":100}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetProductByID(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"测试商品","price":2990,"skus":[{"name":"标准","price":2990,"stock":50}]}`, catID)
	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	c2, w2 := testutil.NewGinContext("GET", fmt.Sprintf("/api/v1/products/%.0f", createResp.Data.ID), "")
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}}
	h.GetByID(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestGetProductByID_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/products/999", "")
	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	h.GetByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListProducts_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/products", "")
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, float64(0), resp.Data["total"])
}

func TestListProducts_ByCategory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"分类内商品","price":1000}`, catID)
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	c2, w2 := testutil.NewGinContext("GET", fmt.Sprintf("/api/v1/products?category_id=%d", catID), "")
	h.List(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["total"])
}

func TestListProducts_FilterByTag(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"热卖商品","price":1000,"tags":"热卖"}`, catID)
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	c2, w2 := testutil.NewGinContext("GET", "/api/v1/products?tag=热卖", "")
	h.List(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["total"])
}

func TestListProducts_SortBySales(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body1 := fmt.Sprintf(`{"category_id":%d,"name":"商品A","price":1000,"sales":10}`, catID)
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body1)
	h.Create(c1)

	body2 := fmt.Sprintf(`{"category_id":%d,"name":"商品B","price":2000,"sales":100}`, catID)
	c2, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body2)
	h.Create(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/products?sort=sales", "")
	h.List(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestSearchProduct(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"有机五常大米","price":4990}`, catID)
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	body2 := fmt.Sprintf(`{"category_id":%d,"name":"有机亚麻籽油","price":5990}`, catID)
	c2, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body2)
	h.Create(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/products/search?keyword=有机", "")
	c3.Set("user_id", uint(1))
	h.Search(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w3.Body.Bytes(), &resp)
	assert.Equal(t, float64(2), resp.Data["total"])
}

func TestSearchProduct_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/products/search?keyword=不存在的商品", "")
	c.Set("user_id", uint(1))
	h.Search(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp.Data["total"])
}

func TestUpdateProductStatus(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"测试商品","price":1000}`, catID)
	c1, w1 := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	c2, w2 := testutil.NewGinContext("PUT", fmt.Sprintf("/api/v1/admin/products/%.0f/status/off", createResp.Data.ID), "")
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", createResp.Data.ID)}, {Key: "status", Value: "off"}}
	h.UpdateStatus(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestHotSearch(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/products/hot-search", "")
	h.HotSearch(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHistory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	catID := createTestCategory(t)

	body := fmt.Sprintf(`{"category_id":%d,"name":"五常大米","price":4990}`, catID)
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/products", body)
	h.Create(c1)

	c2, _ := testutil.NewGinContext("GET", "/api/v1/products/search?keyword=五常大米", "")
	c2.Set("user_id", uint(1))
	h.Search(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/products/search-history", "")
	c3.Set("user_id", uint(1))
	h.SearchHistory(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestClearSearchHistory(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/products/search-history/clear", "")
	c1.Set("user_id", uint(1))
	h.ClearSearchHistory(c1)

	assert.Equal(t, http.StatusOK, w1.Code)
}
