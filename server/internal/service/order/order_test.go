package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/dto/request"
	"gosh/internal/model"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
	}
	cfg := config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"}
	err := database.Init(cfg)
	require.NoError(t, err)
	err = database.DB.AutoMigrate(
		&model.Product{}, &model.ProductSKU{}, &model.Cart{},
		&model.Address{}, &model.Order{}, &model.OrderItem{},
		&model.OrderLog{}, &model.IdempotencyRecord{},
		&model.User{}, &model.PointLog{}, &model.Category{},
	)
	require.NoError(t, err)
}

func seedProduct(t *testing.T, stock int) (productID, skuID uint) {
	t.Helper()
	cat := &model.Category{Name: "测试分类"}
	database.DB.Create(cat)
	p := &model.Product{
		CategoryID: cat.ID,
		Name:       "测试商品",
		Price:      2990,
		Status:     model.ProductStatusOn,
		Images:     []string{"http://example.com/img.jpg"},
	}
	database.DB.Create(p)
	sku := &model.ProductSKU{
		ProductID: p.ID,
		Name:      "标准规格",
		Price:     2990,
		Stock:     stock,
	}
	database.DB.Create(sku)
	return p.ID, sku.ID
}

func seedAddress(t *testing.T, userID uint) {
	t.Helper()
	database.DB.Create(&model.Address{
		UserID:    userID,
		Name:      "张三",
		Phone:     "13800138000",
		Province:  "广东省",
		City:      "深圳市",
		District:  "南山区",
		Detail:    "科技园南区",
		IsDefault: true,
	})
}

func seedCart(t *testing.T, userID, skuID uint, qty int) {
	t.Helper()
	database.DB.Create(&model.Cart{
		UserID:   userID,
		SKUID:    skuID,
		Quantity: qty,
		Selected: true,
	})
}

func seedUserWithPoints(t *testing.T, userID uint, points int) {
	t.Helper()
	database.DB.Create(&model.User{
		Phone:    "13800138000",
		Password: "hashed",
		Role:     model.RoleUser,
		Points:   points,
	})
	database.DB.Model(&model.User{}).Where("id = ?", userID).Update("id", userID)
}

func TestCreate_MissingIdempotentKey(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, err := svc.Create(1, &request.CreateOrderRequest{}, "")
	assert.ErrorIs(t, err, ErrMissingIdempotentKey)
}

func TestCreate_CartEmpty(t *testing.T) {
	setupTestDB(t)
	svc := New()
	seedAddress(t, 1)
	_, err := svc.Create(1, &request.CreateOrderRequest{}, "key-1")
	assert.ErrorIs(t, err, ErrCartEmpty)
}

func TestCreate_NoAddress(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedCart(t, 1, skuID, 1)
	_, err := svc.Create(1, &request.CreateOrderRequest{}, "key-1")
	assert.ErrorIs(t, err, ErrNoDefaultAddress)
}

func TestCreate_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 2)

	resp, err := svc.Create(1, &request.CreateOrderRequest{
		Remark: "测试订单",
	}, "key-create-success")
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusPendingPayment, resp.Status)
	assert.Equal(t, int64(5980), resp.TotalAmount)
	assert.Equal(t, int64(5980), resp.PayAmount)
	assert.Equal(t, "测试订单", resp.Remark)
	require.Len(t, resp.Items, 1)
	assert.Equal(t, skuID, resp.Items[0].SKUID)
	assert.Equal(t, 2, resp.Items[0].Quantity)

	var order model.Order
	database.DB.First(&order, resp.ID)
	assert.Equal(t, model.OrderStatusPendingPayment, order.Status)

	var sku model.ProductSKU
	database.DB.First(&sku, skuID)
	assert.Equal(t, 8, sku.Stock)

	var cartCount int64
	database.DB.Model(&model.Cart{}).Where("user_id = ?", uint(1)).Count(&cartCount)
	assert.Equal(t, int64(0), cartCount)

	var log model.OrderLog
	database.DB.Where("order_id = ?", resp.ID).First(&log)
	assert.Equal(t, "", log.FromStatus)
	assert.Equal(t, model.OrderStatusPendingPayment, log.ToStatus)
}

func TestCreate_BuyNowMode(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)

	_, err := svc.Create(1, &request.CreateOrderRequest{
		Items: []request.CreateOrderItem{{SKUID: skuID, Quantity: 3}},
	}, "key-buy-now")
	require.NoError(t, err)

	var cartCount int64
	database.DB.Model(&model.Cart{}).Where("user_id = ?", uint(1)).Count(&cartCount)
	assert.Equal(t, int64(0), cartCount, "buy-now should not touch cart")
}

func TestCreate_InsufficientStock(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 1)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 5)

	_, err := svc.Create(1, &request.CreateOrderRequest{}, "key-stock")
	assert.ErrorContains(t, err, "insufficient stock")
}

func TestCreate_IdempotentReplay(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)

	first, err := svc.Create(1, &request.CreateOrderRequest{}, "key-replay")
	require.NoError(t, err)

	second, err := svc.Create(1, &request.CreateOrderRequest{}, "key-replay")
	require.NoError(t, err)
	assert.Equal(t, first.ID, second.ID)
	assert.Equal(t, first.OrderNo, second.OrderNo)
	assert.Equal(t, first.TotalAmount, second.TotalAmount)
}

func TestList_Empty(t *testing.T) {
	setupTestDB(t)
	svc := New()
	list, total, err := svc.List(1, &request.ListOrderRequest{Page: 1, Size: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}

func TestList_WithData(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	svc.Create(1, &request.CreateOrderRequest{}, "key-list")

	list, total, err := svc.List(1, &request.ListOrderRequest{Page: 1, Size: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
	assert.Equal(t, model.OrderStatusPendingPayment, list[0].Status)
}

func TestGetByID_Found(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-get")

	resp, err := svc.GetByID(1, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
}

func TestGetByID_NotFound(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, err := svc.GetByID(1, 99999)
	assert.ErrorIs(t, err, ErrOrderNotFound)
}

func TestGetByID_WrongUser(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-wrong-user")

	_, err := svc.GetByID(2, created.ID)
	assert.ErrorIs(t, err, ErrOrderNotBelongToUser)
}

func TestCancel_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 2)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-cancel")

	err := svc.Cancel(1, created.ID, &request.CancelOrderRequest{Reason: "不想要了"})
	require.NoError(t, err)

	var order model.Order
	database.DB.First(&order, created.ID)
	assert.Equal(t, model.OrderStatusCancelled, order.Status)
	assert.Equal(t, "不想要了", order.CancelReason)

	var sku model.ProductSKU
	database.DB.First(&sku, skuID)
	assert.Equal(t, 10, sku.Stock, "stock should be restored")

	var log model.OrderLog
	database.DB.Where("order_id = ? AND to_status = ?", created.ID, model.OrderStatusCancelled).First(&log)
	assert.NotNil(t, log)
}

func TestCancel_WrongStatus(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-cancel-ws")
	svc.Pay(1, created.ID)

	err := svc.Cancel(1, created.ID, &request.CancelOrderRequest{})
	assert.ErrorIs(t, err, ErrInvalidOrderStatus)
}

func TestCancel_WrongUser(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-cancel-wu")

	err := svc.Cancel(2, created.ID, &request.CancelOrderRequest{})
	assert.ErrorIs(t, err, ErrOrderNotBelongToUser)
}

func TestCancel_NotFound(t *testing.T) {
	setupTestDB(t)
	svc := New()
	err := svc.Cancel(1, 99999, &request.CancelOrderRequest{})
	assert.ErrorIs(t, err, ErrOrderNotFound)
}

func TestPay_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-pay")

	err := svc.Pay(1, created.ID)
	require.NoError(t, err)

	var order model.Order
	database.DB.First(&order, created.ID)
	assert.Equal(t, model.OrderStatusPendingDelivery, order.Status)
	assert.NotNil(t, order.PaidAt)
}

func TestPay_WrongStatus(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-pay-ws")
	svc.Pay(1, created.ID)

	err := svc.Pay(1, created.ID)
	assert.ErrorIs(t, err, ErrInvalidOrderStatus)
}

func TestShip_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-ship")
	svc.Pay(1, created.ID)

	err := svc.Ship(created.ID)
	require.NoError(t, err)

	var order model.Order
	database.DB.First(&order, created.ID)
	assert.Equal(t, model.OrderStatusPendingReceipt, order.Status)
	assert.NotNil(t, order.ShippedAt)
}

func TestShip_WrongStatus(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-ship-ws")

	err := svc.Ship(created.ID)
	assert.ErrorIs(t, err, ErrInvalidOrderStatus)
}

func TestShip_NotFound(t *testing.T) {
	setupTestDB(t)
	svc := New()
	err := svc.Ship(99999)
	assert.ErrorIs(t, err, ErrOrderNotFound)
}

func TestConfirm_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-confirm")
	svc.Pay(1, created.ID)
	svc.Ship(created.ID)

	err := svc.Confirm(1, created.ID)
	require.NoError(t, err)

	var order model.Order
	database.DB.First(&order, created.ID)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)
	assert.NotNil(t, order.CompletedAt)
}

func TestConfirm_WrongUser(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-confirm-wu")
	svc.Pay(1, created.ID)
	svc.Ship(created.ID)

	err := svc.Confirm(2, created.ID)
	assert.ErrorIs(t, err, ErrOrderNotBelongToUser)
}

func TestFullLifecycle(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-lifecycle")

	assert.Equal(t, model.OrderStatusPendingPayment, created.Status)
	svc.Pay(1, created.ID)
	svc.Ship(created.ID)
	svc.Confirm(1, created.ID)

	detailed, _ := svc.GetByID(1, created.ID)
	assert.Equal(t, model.OrderStatusCompleted, detailed.Status)
}

func TestRebuy_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{
		Items: []request.CreateOrderItem{{SKUID: skuID, Quantity: 1}},
	}, "key-rebuy")
	svc.Pay(1, created.ID)
	svc.Ship(created.ID)
	svc.Confirm(1, created.ID)

	resp, err := svc.Rebuy(1, created.ID)
	require.NoError(t, err)
	assert.Len(t, resp.Cart.Items, 1)
	assert.Empty(t, resp.SkippedItems)
}

func TestRebuy_WrongStatus(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{
		Items: []request.CreateOrderItem{{SKUID: skuID, Quantity: 1}},
	}, "key-rebuy-ws")

	_, err := svc.Rebuy(1, created.ID)
	assert.ErrorIs(t, err, ErrInvalidOrderStatus)
}

func TestRebuy_SkippedItems(t *testing.T) {
	setupTestDB(t)
	svc := New()
	productID, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	created, _ := svc.Create(1, &request.CreateOrderRequest{
		Items: []request.CreateOrderItem{{SKUID: skuID, Quantity: 1}},
	}, "key-rebuy-sk")
	svc.Pay(1, created.ID)
	svc.Ship(created.ID)
	svc.Confirm(1, created.ID)

	database.DB.Model(&model.Product{}).Where("id = ?", productID).Update("status", model.ProductStatusOff)

	resp, err := svc.Rebuy(1, created.ID)
	require.NoError(t, err)
	assert.Empty(t, resp.Cart.Items)
	assert.Len(t, resp.SkippedItems, 1)
	assert.Equal(t, "商品已下架", resp.SkippedItems[0].Reason)
}

func TestApplyPoints_Success(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	seedUserWithPoints(t, 1, 5000)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-points")

	err := svc.ApplyPoints(1, created.ID, &request.ApplyPointsRequest{Points: 1000})
	require.NoError(t, err)

	var order model.Order
	database.DB.First(&order, created.ID)
	assert.Equal(t, int64(2990-1000), order.PayAmount)
	assert.Equal(t, 1000, order.PointsDeducted)

	var pointLog model.PointLog
	database.DB.Where("order_id = ?", created.ID).First(&pointLog)
	assert.Equal(t, 1000, pointLog.Amount)
}

func TestApplyPoints_Insufficient(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	seedUserWithPoints(t, 1, 100)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-points-ins")

	err := svc.ApplyPoints(1, created.ID, &request.ApplyPointsRequest{Points: 999})
	assert.ErrorIs(t, err, ErrInsufficientPoints)
}

func TestApplyPoints_WrongStatus(t *testing.T) {
	setupTestDB(t)
	svc := New()
	_, skuID := seedProduct(t, 10)
	seedAddress(t, 1)
	seedCart(t, 1, skuID, 1)
	seedUserWithPoints(t, 1, 5000)
	created, _ := svc.Create(1, &request.CreateOrderRequest{}, "key-points-ws")
	svc.Pay(1, created.ID)

	err := svc.ApplyPoints(1, created.ID, &request.ApplyPointsRequest{Points: 100})
	assert.ErrorIs(t, err, ErrInvalidOrderStatus)
}
