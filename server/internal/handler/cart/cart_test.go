package cart

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
	err = database.DB.AutoMigrate(&model.Product{}, &model.ProductSKU{}, &model.Cart{}, &model.Category{})
	require.NoError(t, err)
}

func createTestProduct(t *testing.T) (uint, uint) {
	t.Helper()
	cat := &model.Category{Name: "测试分类"}
	database.DB.Create(cat)
	product := &model.Product{
		CategoryID: cat.ID,
		Name:       "测试商品",
		Price:      2990,
		Status:     model.ProductStatusOn,
	}
	database.DB.Create(product)
	sku := &model.ProductSKU{
		ProductID: product.ID,
		Name:      "标准规格",
		Price:     2990,
		Stock:     100,
	}
	database.DB.Create(sku)
	return product.ID, sku.ID
}

func TestAddCart_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":2}`, skuID))
	c.Set("user_id", uint(1))
	h.Add(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, float64(skuID), resp.Data["sku_id"])
	assert.Equal(t, float64(2), resp.Data["quantity"])
}

func TestAddCart_DuplicateMerge(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, _ := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":2}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":3}`, skuID))
	c2.Set("user_id", uint(1))
	h.Add(c2)

	assert.Equal(t, http.StatusCreated, w2.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(5), resp.Data["quantity"])
}

func TestAddCart_SKUNotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/cart", `{"sku_id":999,"quantity":1}`)
	c.Set("user_id", uint(1))
	h.Add(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddCart_MissingParams(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/cart", `{"quantity":1}`)
	c.Set("user_id", uint(1))
	h.Add(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListCart_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/cart", "")
	c.Set("user_id", uint(1))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp.Data["items"].([]interface{})
	assert.Equal(t, 0, len(items))
}

func TestListCart_WithItems(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, _ := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":1}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("GET", "/api/v1/cart", "")
	c2.Set("user_id", uint(1))
	h.List(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	items := resp.Data["items"].([]interface{})
	assert.Equal(t, 1, len(items))
}

func TestUpdateCart_Quantity(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":1}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	var createResp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	cartID := createResp.Data["id"]

	c2, w2 := testutil.NewGinContext("PUT", fmt.Sprintf("/api/v1/cart/%.0f", cartID), `{"quantity":5}`)
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", cartID)}}
	h.Update(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(5), resp.Data["quantity"])
}

func TestDeleteCart(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":1}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	var createResp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	cartID := createResp.Data["id"]

	c2, w2 := testutil.NewGinContext("DELETE", fmt.Sprintf("/api/v1/cart/%.0f", cartID), "")
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%.0f", cartID)}}
	h.Delete(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestDeleteCart_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("DELETE", "/api/v1/cart/999", "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCartCount(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, _ := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":2}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("GET", "/api/v1/cart/count", "")
	c2.Set("user_id", uint(1))
	h.Count(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["count"])
}

func TestMergeCart(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID1 := createTestProduct(t)
	_, skuID2 := createTestProduct(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/cart/merge", fmt.Sprintf(
		`{"items":[{"sku_id":%d,"quantity":2},{"sku_id":%d,"quantity":3}]}`, skuID1, skuID2))
	c.Set("user_id", uint(1))
	h.Merge(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp.Data["items"].([]interface{})
	assert.Equal(t, 2, len(items))
}

func TestSelectCart(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProduct(t)

	c1, _ := testutil.NewGinContext("POST", "/api/v1/cart", fmt.Sprintf(`{"sku_id":%d,"quantity":1}`, skuID))
	c1.Set("user_id", uint(1))
	h.Add(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/cart/select", fmt.Sprintf(`{"sku_ids":[%d],"selected":false}`, skuID))
	c2.Set("user_id", uint(1))
	h.Select(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
}
