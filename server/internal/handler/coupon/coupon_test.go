package coupon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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
	err = database.DB.AutoMigrate(&model.Coupon{}, &model.UserCoupon{})
	require.NoError(t, err)
}

func TestCreateCoupon_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	startAt := time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")
	endAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	body := fmt.Sprintf(`{"name":"满100减20","type":"full_reduce","condition":10000,"discount":2000,"total_count":100,"per_limit":1,"start_at":"%s","end_at":"%s"}`, startAt, endAt)

	c, w := testutil.NewGinContext("POST", "/api/v1/admin/coupons", body)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "满100减20", resp.Data["name"])
	assert.Equal(t, float64(100), resp.Data["remain_count"])
}

func TestCreateCoupon_InvalidTime(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"name":"满100减20","type":"full_reduce","condition":10000,"discount":2000,"start_at":"invalid","end_at":"invalid"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/coupons", body)
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func createTestCoupon(t *testing.T, remain int) *model.Coupon {
	t.Helper()
	coupon := &model.Coupon{
		Name:        "满100减20",
		Type:        model.CouponTypeFullReduce,
		Condition:   10000,
		Discount:    2000,
		TotalCount:  100,
		RemainCount: remain,
		PerLimit:    1,
		StartAt:     time.Now().Add(-time.Hour),
		EndAt:       time.Now().Add(24 * time.Hour),
		Status:      model.CouponStatusActive,
	}
	database.DB.Create(coupon)
	return coupon
}

func TestReceiveCoupon_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := createTestCoupon(t, 10)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/coupons/%d/receive", coupon.ID), "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", coupon.ID)}}
	h.Receive(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var updated model.Coupon
	database.DB.First(&updated, coupon.ID)
	assert.Equal(t, 9, updated.RemainCount)
}

func TestReceiveCoupon_SoldOut(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := createTestCoupon(t, 0)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/coupons/%d/receive", coupon.ID), "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", coupon.ID)}}
	h.Receive(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReceiveCoupon_Duplicate(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := createTestCoupon(t, 10)

	c1, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/coupons/%d/receive", coupon.ID), "")
	c1.Set("user_id", uint(1))
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", coupon.ID)}}
	h.Receive(c1)

	c2, w2 := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/coupons/%d/receive", coupon.ID), "")
	c2.Set("user_id", uint(1))
	c2.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", coupon.ID)}}
	h.Receive(c2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestReceiveCoupon_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("POST", "/api/v1/coupons/999/receive", "")
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	h.Receive(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCalculate_FullReduce(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := createTestCoupon(t, 10)

	body := fmt.Sprintf(`{"order_amount":20000,"coupon_id":%d}`, coupon.ID)
	c, w := testutil.NewGinContext("POST", "/api/v1/coupons/calculate", body)
	h.Calculate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(20000), resp.Data["original_amount"])
	assert.Equal(t, float64(2000), resp.Data["discount_amount"])
	assert.Equal(t, float64(18000), resp.Data["pay_amount"])
}

func TestCalculate_NotApplicable(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := createTestCoupon(t, 10)

	body := fmt.Sprintf(`{"order_amount":5000,"coupon_id":%d}`, coupon.ID)
	c, w := testutil.NewGinContext("POST", "/api/v1/coupons/calculate", body)
	h.Calculate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculate_Discount(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	coupon := &model.Coupon{
		Name:       "9折券",
		Type:       model.CouponTypeDiscount,
		Condition:  0,
		Discount:   90,
		TotalCount: 100,
		RemainCount: 100,
		PerLimit:   1,
		StartAt:    time.Now().Add(-time.Hour),
		EndAt:      time.Now().Add(24 * time.Hour),
		Status:     model.CouponStatusActive,
	}
	database.DB.Create(coupon)

	body := fmt.Sprintf(`{"order_amount":20000,"coupon_id":%d}`, coupon.ID)
	c, w := testutil.NewGinContext("POST", "/api/v1/coupons/calculate", body)
	h.Calculate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(2000), resp.Data["discount_amount"])
	assert.Equal(t, float64(18000), resp.Data["pay_amount"])
}
