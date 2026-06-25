package order

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
	err = database.DB.AutoMigrate(
		&model.Product{}, &model.ProductSKU{}, &model.Cart{},
		&model.Address{}, &model.Order{}, &model.OrderItem{},
		&model.OrderLog{}, &model.IdempotencyRecord{}, &model.Category{},
		&model.User{}, &model.PointLog{},
	)
	require.NoError(t, err)
}

func createTestProductWithStock(t *testing.T, stock int) (uint, uint) {
	t.Helper()
	cat := &model.Category{Name: "测试分类"}
	database.DB.Create(cat)
	product := &model.Product{
		CategoryID: cat.ID,
		Name:       "测试商品",
		Price:      2990,
		Status:     model.ProductStatusOn,
		Images:     []string{"http://example.com/img.jpg"},
	}
	database.DB.Create(product)
	sku := &model.ProductSKU{
		ProductID: product.ID,
		Name:      "标准规格",
		Price:     2990,
		Stock:     stock,
	}
	database.DB.Create(sku)
	return product.ID, sku.ID
}

func addToCart(t *testing.T, userID, skuID uint, qty int) {
	t.Helper()
	cart := &model.Cart{UserID: userID, SKUID: skuID, Quantity: qty, Selected: true}
	database.DB.Create(cart)
}

func createDefaultAddress(t *testing.T, userID uint) uint {
	t.Helper()
	addr := &model.Address{
		UserID:    userID,
		Name:      "张三",
		Phone:     "13800138000",
		Province:  "广东省",
		City:      "深圳市",
		District:  "南山区",
		Detail:    "科技园南区",
		IsDefault: true,
	}
	database.DB.Create(addr)
	return addr.ID
}

func TestCreateOrder_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProductWithStock(t, 100)
	createDefaultAddress(t, 1)
	addToCart(t, 1, skuID, 2)

	c, w := testutil.NewGinContext("POST", "/api/v1/orders", `{"remark":"请尽快发货"}`)
	c.Request.Header.Set("Idempotent-Key", "test-key-1")
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "pending_payment", resp.Data["status"])
	assert.Equal(t, float64(5980), resp.Data["total_amount"])
	assert.NotEmpty(t, resp.Data["order_no"])
}

func TestCreateOrder_MissingIdempotentKey(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProductWithStock(t, 100)
	createDefaultAddress(t, 1)
	addToCart(t, 1, skuID, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_EmptyCart(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	createDefaultAddress(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c.Request.Header.Set("Idempotent-Key", "test-key-2")
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_NoAddress(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProductWithStock(t, 100)
	addToCart(t, 1, skuID, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c.Request.Header.Set("Idempotent-Key", "test-key-3")
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProductWithStock(t, 1)
	createDefaultAddress(t, 1)
	addToCart(t, 1, skuID, 5)

	c, w := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c.Request.Header.Set("Idempotent-Key", "test-key-4")
	c.Set("user_id", uint(1))
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_Idempotency(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	_, skuID := createTestProductWithStock(t, 100)
	createDefaultAddress(t, 1)
	addToCart(t, 1, skuID, 1)

	c1, w1 := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c1.Request.Header.Set("Idempotent-Key", "same-key")
	c1.Set("user_id", uint(1))
	h.Create(c1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/orders", `{}`)
	c2.Request.Header.Set("Idempotent-Key", "same-key")
	c2.Set("user_id", uint(1))
	h.Create(c2)

	assert.Equal(t, http.StatusCreated, w2.Code)
}

func TestListOrders_Empty(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/orders", "")
	c.Set("user_id", uint(1))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp.Data["total"])
}

func createTestOrder(t *testing.T, userID uint) uint {
	t.Helper()
	_, skuID := createTestProductWithStock(t, 100)
	createDefaultAddress(t, userID)
	addToCart(t, userID, skuID, 2)

	order := &model.Order{
		OrderNo:       fmt.Sprintf("TEST%d", userID),
		UserID:        userID,
		Status:        model.OrderStatusPendingPayment,
		TotalAmount:   5980,
		PayAmount:     5980,
		AddressName:   "张三",
		AddressPhone:  "13800138000",
		AddressDetail: `{"name":"张三","phone":"13800138000","province":"广东省","city":"深圳市","district":"南山区","detail":"科技园南区"}`,
	}
	database.DB.Create(order)
	database.DB.Create(&model.OrderItem{
		OrderID:     order.ID,
		SKUID:       skuID,
		ProductName: "测试商品",
		SKUName:     "标准规格",
		Price:       2990,
		Quantity:    2,
		Subtotal:    5980,
	})
	database.DB.Create(&model.OrderLog{
		OrderID:    order.ID,
		FromStatus: "",
		ToStatus:   model.OrderStatusPendingPayment,
		Operator:   "user:1",
		Note:       "订单创建",
	})
	return order.ID
}

func TestGetOrderByID(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("GET", fmt.Sprintf("/api/v1/orders/%d", orderID), "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "pending_payment", resp.Data["status"])
}

func TestGetOrderByID_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/orders/999", "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	h.GetByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCancelOrder(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{"reason":"不想要了"}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var order model.Order
	database.DB.First(&order, orderID)
	assert.Equal(t, model.OrderStatusCancelled, order.Status)
}

func TestCancelOrder_WrongStatus(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	// First cancel
	c1, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{}`)
	c1.Set("user_id", uint(1))
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c1)

	// Second cancel should fail
	c2, w2 := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{}`)
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestPayOrder(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/pay", orderID), `{}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Pay(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFullOrderFlow(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	// Pay
	c1, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/pay", orderID), `{}`)
	c1.Set("user_id", uint(1))
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Pay(c1)

	// Ship (admin)
	c2, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/ship", orderID), `{}`)
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Ship(c2)

	// Confirm receipt
	c3, w3 := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/confirm", orderID), `{}`)
	c3.Set("user_id", uint(1))
	c3.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Confirm(c3)

	assert.Equal(t, http.StatusOK, w3.Code)

	var order model.Order
	database.DB.First(&order, orderID)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)
}

func TestApplyPoints_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	database.DB.Create(&model.User{Phone: "13800000001", Password: "test"})
	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 1000)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/apply-points", orderID),
		`{"points":500}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.ApplyPoints(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var order model.Order
	database.DB.First(&order, orderID)
	assert.Equal(t, 500, order.PointsDeducted)
	assert.Equal(t, int64(5980-500), order.PayAmount)

	var user model.User
	database.DB.Select("points").First(&user, 1)
	assert.Equal(t, 500, user.Points)

	var log model.PointLog
	err := database.DB.Where("user_id = ? AND type = ?", 1, model.PointTypeSpend).First(&log).Error
	assert.NoError(t, err)
	assert.Equal(t, 500, log.Amount)
	assert.Equal(t, 500, log.Balance)
}

func TestApplyPoints_InsufficientPoints(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	database.DB.Create(&model.User{Phone: "13800000002", Password: "test"})
	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 10)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/apply-points", orderID),
		`{"points":500}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.ApplyPoints(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestApplyPoints_WrongOrderStatus(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	database.DB.Create(&model.User{Phone: "13800000003", Password: "test"})
	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 1000)

	c1, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{}`)
	c1.Set("user_id", uint(1))
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c1)

	c2, w2 := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/apply-points", orderID),
		`{"points":100}`)
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.ApplyPoints(c2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestApplyPoints_NotExceedPayAmount(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	database.DB.Create(&model.User{Phone: "13800000004", Password: "test"})
	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 99999)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/apply-points", orderID),
		`{"points":99999}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.ApplyPoints(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var order model.Order
	database.DB.First(&order, orderID)
	assert.True(t, order.PayAmount >= 0)
}

func TestOrderListFilterByStatus(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("GET", "/api/v1/orders?status=pending_payment", "")
	c.Set("user_id", uint(1))
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["total"])

	// Cancel and check filter
	c2, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{}`)
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c2)

	c3, w3 := testutil.NewGinContext("GET", "/api/v1/orders?status=cancelled", "")
	c3.Set("user_id", uint(1))
	h.List(c3)

	assert.Equal(t, http.StatusOK, w3.Code)
	json.Unmarshal(w3.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.Data["total"])
}
