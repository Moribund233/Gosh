package payment

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
		&model.User{},
		&model.Product{}, &model.ProductSKU{},
		&model.Address{}, &model.Order{}, &model.OrderItem{},
		&model.OrderLog{}, &model.IdempotencyRecord{}, &model.Cart{},
		&model.Payment{}, &model.PointLog{},
	)
	require.NoError(t, err)
}

func createTestUser(t *testing.T, userID uint) {
	t.Helper()
	database.DB.Create(&model.User{
		Phone:    fmt.Sprintf("138001380%02d", userID),
		Password: "hash",
		Nickname: fmt.Sprintf("User%d", userID),
	})
}

func createTestOrder(t *testing.T, userID uint) *model.Order {
	t.Helper()
	createTestUser(t, userID)
	cat := &model.Category{Name: "测试分类"}
	database.DB.Create(cat)
	product := &model.Product{
		CategoryID: cat.ID,
		Name:       "测试商品",
		Price:      2990,
		Status:     model.ProductStatusOn,
		Images:     "http://example.com/img.jpg",
	}
	database.DB.Create(product)
	sku := &model.ProductSKU{
		ProductID: product.ID,
		Name:      "标准规格",
		Price:     2990,
		Stock:     100,
	}
	database.DB.Create(sku)

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

	cart := &model.Cart{UserID: userID, SKUID: sku.ID, Quantity: 2, Selected: true}
	database.DB.Create(cart)

	orderNo := fmt.Sprintf("TEST%04d", userID)
	order := &model.Order{
		OrderNo:    orderNo,
		UserID:     userID,
		Status:     model.OrderStatusPendingPayment,
		TotalAmount: 5980,
		PayAmount:   5980,
		AddressName: "张三",
		AddressPhone: "13800138000",
	}
	database.DB.Create(order)
	return order
}

func TestGetMethods(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/payment/methods", "")
	h.GetMethods(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.Data, 3)
	assert.Equal(t, "mock", resp.Data[0]["method"])
	assert.Equal(t, "wechat", resp.Data[1]["method"])
	assert.Equal(t, "alipay", resp.Data[2]["method"])
}

func TestPay_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	c.Set("user_id", uint(1))
	h.Pay(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Data["transaction_no"])
	assert.Equal(t, "success", resp.Data["status"])
	assert.Equal(t, float64(5980), resp.Data["pay_amount"])

	var updated model.Order
	database.DB.First(&updated, order.ID)
	assert.Equal(t, model.OrderStatusPendingDelivery, updated.Status)
}

func TestPay_OrderNotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		`{"order_no":"NONEXISTENT","method":"mock"}`)
	c.Set("user_id", uint(1))
	h.Pay(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPay_OrderNotBelongToUser(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	c.Set("user_id", uint(2))
	h.Pay(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPay_InvalidOrderStatus(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	firstPay, _ := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	firstPay.Set("user_id", uint(1))
	h.Pay(firstPay)

	secondPay, w2 := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	secondPay.Set("user_id", uint(1))
	h.Pay(secondPay)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestPay_InvalidMethod(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"invalid"}`, order.OrderNo))
	c.Set("user_id", uint(1))
	h.Pay(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStatus_Found(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	payC, _ := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	payC.Set("user_id", uint(1))
	h.Pay(payC)

	c, w := testutil.NewGinContext("GET", fmt.Sprintf("/api/v1/payment/status/%s", order.OrderNo), "")
	c.Params = []gin.Param{{Key: "order_no", Value: order.OrderNo}}
	h.GetStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "success", resp.Data["status"])
}

func TestGetStatus_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/payment/status/NONEXISTENT", "")
	c.Params = []gin.Param{{Key: "order_no", Value: "NONEXISTENT"}}
	h.GetStatus(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCallback_Mock_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	body := fmt.Sprintf(`{"transaction_no":"MOCKTEST001","order_no":"%s","amount":5980,"timestamp":1700000000,"sign":"test","method":"mock"}`, order.OrderNo)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/callback/mock", body)
	c.Params = []gin.Param{{Key: "method", Value: "mock"}}
	h.Callback(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPay_AuditLog(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	c.Set("user_id", uint(1))
	h.Pay(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var log model.OrderLog
	err := database.DB.Where("order_id = ? AND to_status = ?", order.ID, model.OrderStatusPendingDelivery).First(&log).Error
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStatusPendingPayment, log.FromStatus)
	assert.Contains(t, log.Note, "mock")
}

func TestRefund_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	payC, _ := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	payC.Set("user_id", uint(1))
	h.Pay(payC)

	var payment model.Payment
	database.DB.Where("order_no = ?", order.OrderNo).First(&payment)

	c, w := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		fmt.Sprintf(`{"transaction_no":"%s","amount":5980,"reason":"商品质量问题"}`, payment.TransactionNo))
	c.Set("user_id", uint(1))
	h.Refund(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedPayment model.Payment
	database.DB.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

	var updatedOrder model.Order
	database.DB.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusCancelled, updatedOrder.Status)
}

func TestRefund_InvalidTransactionNo(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		`{"transaction_no":"NONEXISTENT","amount":100,"reason":"test"}`)
	c.Set("user_id", uint(1))
	h.Refund(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRefund_DuplicateRefund(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	payC, _ := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	payC.Set("user_id", uint(1))
	h.Pay(payC)

	var payment model.Payment
	database.DB.Where("order_no = ?", order.OrderNo).First(&payment)

	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		fmt.Sprintf(`{"transaction_no":"%s","amount":5980,"reason":"first"}`, payment.TransactionNo))
	c1.Set("user_id", uint(1))
	h.Refund(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		fmt.Sprintf(`{"transaction_no":"%s","amount":5980,"reason":"second"}`, payment.TransactionNo))
	c2.Set("user_id", uint(1))
	h.Refund(c2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestPay_AmountInCents(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	c, w := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	c.Set("user_id", uint(1))
	h.Pay(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(5980), resp.Data["pay_amount"])
}
