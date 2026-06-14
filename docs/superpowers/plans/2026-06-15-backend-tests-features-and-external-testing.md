# Backend Tests, Features & External Testing Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete remaining backend features (refund, point deduction), add missing tests (upload handler), and establish external testing infrastructure (k6/Python).

**Architecture:** All new features follow existing Handler → Service(interface) → Repository(interface) layered pattern. Tests use in-memory SQLite + `testutil.NewGinContext`. External testing runs against the real server with SQLite.

**Tech Stack:** Go 1.x, Gin, GORM, testify, k6 (or Python + requests)

---

### Task 1: Upload Handler Tests

**Files:**
- Create: `server/internal/handler/upload/upload_test.go`

- [ ] **Step 1: Write tests for Upload (multipart form)**

```go
package upload

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gosh/internal/config"
	"gosh/internal/testutil"
)

func setupUploadTest(t *testing.T) {
	t.Helper()
	config.AppConfig = &config.Config{
		Upload: config.UploadConfig{
			Dir:     filepath.Join(os.TempDir(), "gosh-test-uploads"),
			MaxSize: 10,
		},
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}
	os.MkdirAll(config.AppConfig.Upload.Dir, 0755)
}

func TestUpload_Success(t *testing.T) {
	setupUploadTest(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write([]byte("fake-image-data"))
	writer.Close()

	c, w := testutil.NewGinContext("POST", "/api/v1/upload", body.String())
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	Upload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Data["url"], "/uploads/")
	assert.NotNil(t, resp.Data["size"])
}

func TestUpload_MissingFile(t *testing.T) {
	setupUploadTest(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/upload", "")
	Upload(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpload_UnsupportedType(t *testing.T) {
	setupUploadTest(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.exe")
	part.Write([]byte("bad"))
	writer.Close()

	c, w := testutil.NewGinContext("POST", "/api/v1/upload", body.String())
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	Upload(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpload_FileTooLarge(t *testing.T) {
	config.AppConfig = &config.Config{
		Upload: config.UploadConfig{Dir: os.TempDir(), MaxSize: 1},
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "large.jpg")
	largeData := make([]byte, 2*1024*1024) // 2MB > 1MB limit
	part.Write(largeData)
	writer.Close()

	c, w := testutil.NewGinContext("POST", "/api/v1/upload", body.String())
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	Upload(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadBase64_Success(t *testing.T) {
	setupUploadTest(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/upload/base64", `{"data":"data:image/png;base64,`+
		`iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGA"}`)
	UploadBase64(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Data["url"], "/uploads/")
}

func TestUploadBase64_MissingData(t *testing.T) {
	setupUploadTest(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/upload/base64", `{}`)
	UploadBase64(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadBase64_InvalidData(t *testing.T) {
	setupUploadTest(t)

	c, w := testutil.NewGinContext("POST", "/api/v1/upload/base64", `{"data":"not-valid-base64!!!"}`)
	UploadBase64(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

- [ ] **Step 2: Run upload handler tests to verify they pass**

Run: `cd server && go test ./internal/handler/upload/ -v`
Expected: ALL PASS

- [ ] **Step 3: Run full test suite to confirm no regressions**

Run: `cd server && go test ./...`
Expected: ALL PASS

---

### Task 2: Refund Feature Implementation

**Files:**
- Modify: `server/internal/service/payment/payment.go` — implement Refund
- Modify: `server/internal/handler/payment/payment.go` — add Refund handler
- Modify: `server/internal/router/router.go` — register refund route
- Create: `server/internal/dto/request/payment.go` — already exists, verify RefundRequest
- Create: `server/internal/dto/response/payment.go` — add RefundResponse (optional)
- Create: `server/internal/handler/payment/payment_test.go` — add refund tests

**Design:**
- Refund only allowed for paid orders (status: pending_delivery or later, but not cancelled)
- Validate transaction_no exists and belongs to the user (or admin)
- Create a refund payment record (status=refunded)
- Update order status to reflect refund
- Write OrderLog audit entry
- Only super_admin/support can process refund (per roadmap)

- [ ] **Step 1: Add refund-related errors and model constants**

Add to `server/internal/model/payment.go`:
```go
const (
	PaymentStatusRefunded = "refunded"
)
```

- [ ] **Step 2: Add Refund method to payment Service interface and implementation**

Modify `server/internal/service/payment/payment.go`:

```go
var (
	// existing errors...
	ErrRefundNotImplemented = errors.New("refund not yet implemented") // keep for backward compat
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrRefundExists         = errors.New("refund already processed for this payment")
)

type Service interface {
	GetMethods() []response.PaymentMethodResponse
	Pay(userID uint, req *request.PayRequest) (*response.PaymentResponse, error)
	ProcessCallback(method string, body []byte) error
	GetStatus(orderNo string) (*response.PaymentResponse, error)
	Refund(userID uint, req *request.RefundRequest) error
}
```

Add implementation after GetStatus:
```go
func (s *service) Refund(userID uint, req *request.RefundRequest) error {
	payment, err := s.paymentRepo.FindByTransactionNo(req.TransactionNo)
	if err != nil {
		return ErrPaymentNotFound
	}
	if payment.Status == model.PaymentStatusRefunded {
		return ErrRefundExists
	}
	if payment.Status != model.PaymentStatusSuccess {
		return fmt.Errorf("payment not in success status")
	}

	order, err := s.orderRepo.FindByOrderNo(payment.OrderNo)
	if err != nil {
		return ErrOrderNotFound
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 标记支付记录为已退款
		if err := tx.Model(&model.Payment{}).
			Where("id = ?", payment.ID).
			Update("status", model.PaymentStatusRefunded).Error; err != nil {
			return err
		}

		// 恢复库存
		var items []model.OrderItem
		if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Model(&model.ProductSKU{}).
				Where("id = ?", item.SKUID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// 更新订单状态
		res := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"status":        model.OrderStatusCancelled,
				"cancelled_at":  time.Now(),
				"cancel_reason": fmt.Sprintf("退款: %s", req.Reason),
				"version":       order.Version + 1,
			})
		if res.RowsAffected == 0 {
			return ErrInvalidOrderStatus
		}

		return tx.Create(&model.OrderLog{
			OrderID:    order.ID,
			FromStatus: order.Status,
			ToStatus:   model.OrderStatusCancelled,
			Operator:   fmt.Sprintf("admin:%d", userID),
			Note:       fmt.Sprintf("退款处理: %s (交易号: %s)", req.Reason, req.TransactionNo),
		}).Error
	})
}
```

- [ ] **Step 3: Add Refund handler**

Add to `server/internal/handler/payment/payment.go`:

```go
func (h *Handler) Refund(c *gin.Context) {
	var req request.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	err := h.svc.Refund(userID.(uint), &req)
	if err != nil {
		if err == svc.ErrPaymentNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrRefundExists {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "refund failed")
		return
	}
	response.Success(c, nil)
}
```

- [ ] **Step 4: Register refund route in router**

Add to `server/internal/router/router.go` inside the `admin` block (after coupon create):
```go
// Payment management
admin.POST("/payment/refund", paymentH.Refund)
```

- [ ] **Step 5: Write refund tests**

Add to `server/internal/handler/payment/payment_test.go`:

```go
func TestRefund_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	order := createTestOrder(t, 1)

	// Pay first
	payC, _ := testutil.NewGinContext("POST", "/api/v1/payment/pay",
		fmt.Sprintf(`{"order_no":"%s","method":"mock"}`, order.OrderNo))
	payC.Set("user_id", uint(1))
	h.Pay(payC)

	var payment model.Payment
	database.DB.Where("order_no = ?", order.OrderNo).First(&payment)

	// Refund
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

	// First refund
	c1, _ := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		fmt.Sprintf(`{"transaction_no":"%s","amount":5980,"reason":"test"}`, payment.TransactionNo))
	c1.Set("user_id", uint(1))
	h.Refund(c1)

	// Second refund should fail
	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/payment/refund",
		fmt.Sprintf(`{"transaction_no":"%s","amount":5980,"reason":"test again"}`, payment.TransactionNo))
	c2.Set("user_id", uint(1))
	h.Refund(c2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
```

(No new imports needed — `testutil.NewGinContext` already returns `*gin.Context, *httptest.ResponseRecorder`)

- [ ] **Step 6: Run tests to verify refund works**

Run: `cd server && go test ./internal/handler/payment/ -v`
Expected: ALL PASS

- [ ] **Step 7: Run full test suite**

Run: `cd server && go test ./...`
Expected: ALL PASS

---

### Task 3: Point Deduction at Checkout

**Files:**
- Modify: `server/internal/model/order.go` — add PointsDeducted field
- Modify: `server/internal/service/order/order.go` — add ApplyPoints method
- Modify: `server/internal/handler/order/order.go` — add ApplyPoints handler
- Modify: `server/internal/router/router.go` — register route
- Create: `server/internal/dto/request/order.go` — add ApplyPointsRequest (or use existing)
- Modify: `server/internal/handler/order/order_test.go` — add tests

**Design:**
- 1 point = 1分 (1 cent) monetary value
- Rate: earn 1 point per 100分 spent (1% cashback); deduction 1 point = 1分
- Apply points before payment, only for pending_payment orders
- Max deduction = min(user_points, pay_amount)
- Updates order.PointsDeducted, order.DiscountAmount, order.PayAmount

- [ ] **Step 1: Add PointsDeducted to Order model**

Modify `server/internal/model/order.go` — add field after `CancelReason`:
```go
PointsDeducted int  `gorm:"default:0" json:"points_deducted"`
```

- [ ] **Step 2: Add ApplyPointsRequest DTO**

Add to `server/internal/dto/request/order.go`:
```go
type ApplyPointsRequest struct {
	Points int `json:"points" binding:"required,min=1"`
}
```

- [ ] **Step 3: Add ApplyPoints to order Service interface**

Add to existing Service interface in `server/internal/service/order/order.go`:
```go
ApplyPoints(userID, orderID uint, req *request.ApplyPointsRequest) error
```

Add error:
```go
var (
	ErrInsufficientPoints = errors.New("insufficient points")
)
```

Add implementation:
```go
func (s *service) ApplyPoints(userID, orderID uint, req *request.ApplyPointsRequest) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusPendingPayment {
		return ErrInvalidOrderStatus
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 检查用户积分余额
		var user model.User
		if err := tx.Select("points").First(&user, userID).Error; err != nil {
			return err
		}
		if user.Points < req.Points {
			return ErrInsufficientPoints
		}

		// 计算可抵扣金额（1 point = 1分）
		pointsAmount := int64(req.Points)
		if pointsAmount > order.PayAmount {
			pointsAmount = order.PayAmount
		}

		// 扣除用户积分
		res := tx.Model(&model.User{}).
			Where("id = ? AND points >= ?", userID, req.Points).
			Update("points", gorm.Expr("points - ?", req.Points))
		if res.RowsAffected == 0 {
			return ErrInsufficientPoints
		}

		// 更新订单折扣和实付金额
		newDiscount := order.DiscountAmount + pointsAmount
		newPayAmount := order.TotalAmount + order.ShippingFee - newDiscount
		if newPayAmount < 0 {
			newPayAmount = 0
		}

		if err := tx.Model(&model.Order{}).
			Where("id = ? AND version = ?", order.ID, order.Version).
			Updates(map[string]interface{}{
				"discount_amount": newDiscount,
				"pay_amount":      newPayAmount,
				"points_deducted": order.PointsDeducted + req.Points,
				"version":         order.Version + 1,
			}).Error; err != nil {
			return err
		}

		// 记录积分流水
		var updatedUser model.User
		tx.Select("points").First(&updatedUser, userID)
		return tx.Create(&model.PointLog{
			UserID:  userID,
			Type:    model.PointTypeSpend,
			Amount:  req.Points,
			Balance: updatedUser.Points,
			OrderID: &orderID,
			Note:    fmt.Sprintf("订单抵扣 %d 积分，减免 %.2f 元", req.Points, float64(pointsAmount)/100),
		}).Error
	})
}
```

- [ ] **Step 4: Add ApplyPoints handler**

Add to `server/internal/handler/order/order.go`:

```go
func (h *Handler) ApplyPoints(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	var req request.ApplyPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.ApplyPoints(userID.(uint), uint(id), &req); err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus || err == svc.ErrInsufficientPoints {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "apply points failed")
		return
	}
	response.Success(c, nil)
}
```

- [ ] **Step 5: Register ApplyPoints route**

Add to `server/internal/router/router.go` inside the auth group after orders:
```go
auth.POST("/orders/:id/apply-points", orderH.ApplyPoints)
```

- [ ] **Step 6: Write ApplyPoints tests**

Add to `server/internal/handler/order/order_test.go`:

```go
func TestApplyPoints_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()
	orderID := createTestOrder(t, 1)

	// Give user some points
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

	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 1000)

	// Cancel the order first
	c1, _ := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/cancel", orderID), `{}`)
	c1.Set("user_id", uint(1))
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.Cancel(c1)

	// Try to apply points
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

	database.DB.Model(&model.User{}).Where("id = ?", 1).Update("points", 99999)

	c, w := testutil.NewGinContext("POST", fmt.Sprintf("/api/v1/orders/%d/apply-points", orderID),
		`{"points":99999}`)
	c.Set("user_id", uint(1))
	c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", orderID)}}
	h.ApplyPoints(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var order model.Order
	database.DB.First(&order, orderID)
	// PayAmount should not go below 0
	assert.True(t, order.PayAmount >= 0)
}
```

- [ ] **Step 7: Run tests to verify ApplyPoints works**

Run: `cd server && go test ./internal/handler/order/ -v`
Expected: ALL PASS

- [ ] **Step 8: Run full test suite**

Run: `cd server && go test ./...`
Expected: ALL PASS

---

### Task 4: Docker Compose for PostgreSQL (Load Testing)

**Files:**
- Create: `server/docker-compose.yml` — PostgreSQL for load test environment

**Rationale:** SQLite doesn't handle concurrent writes well. For k6 load testing, switch to PostgreSQL via Docker Compose.

- [ ] **Step 1: Create docker-compose.yml**

Create `server/docker-compose.yml`:
```yaml
version: "3.9"

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: gosh
      POSTGRES_PASSWORD: gosh
      POSTGRES_DB: gosh
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

- [ ] **Step 2: Start PostgreSQL**

Run: `cd server && docker compose up -d`
Verify: `docker compose ps`

---

### Task 5: External Testing Scripts

**Files:**
- Create: `server/scripts/api_test.py` — Python end-to-end API test (runs against SQLite dev)
- Create: `server/scripts/k6_test.js` — k6 load test script (runs against PostgreSQL)

**Note:** Python smoke test uses dev config (SQLite, default). K6 load test uses prod config (PostgreSQL).

- [ ] **Step 1: Create Python end-to-end test script**

Create `server/scripts/api_test.py`:

```python
#!/usr/bin/env python3
"""
Gosh Mall API 端到端测试脚本

用法:
  pip install requests
  python scripts/api_test.py

环境变量:
  BASE_URL (默认: http://localhost:8080/api/v1)
"""

import os
import sys
import json
import time
import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8080/api/v1")

def log(step, status, detail=""):
    icon = "✅" if status == "PASS" else "❌"
    print(f"{icon} [{status}] {step}  {detail}")

def check(label, ok, detail=""):
    if ok:
        log(label, "PASS", detail)
    else:
        log(label, "FAIL", detail)
        sys.exit(1)

def main():
    session = requests.Session()

    # === 1. Health Check ===
    r = session.get(f"{BASE_URL.replace('/api/v1', '')}/health")
    check("Health check", r.status_code == 200)

    # === 2. Register ===
    phone = f"138{int(time.time()) % 10000000000:010d}"[-11:]
    r = session.post(f"{BASE_URL}/user/register", json={
        "phone": phone,
        "password": "test1234",
        "nickname": "测试用户"
    })
    check("User register", r.status_code == 201, f"phone={phone}")
    token = r.json().get("data", {}).get("token")
    check("Token returned", bool(token))

    session.headers["Authorization"] = f"Bearer {token}"

    # === 3. Get Profile ===
    r = session.get(f"{BASE_URL}/user/profile")
    check("Get profile", r.status_code == 200)
    check("Profile has nickname", r.json()["data"]["nickname"] == "测试用户")

    # === 4. Create Address ===
    r = session.post(f"{BASE_URL}/addresses", json={
        "name": "张三", "phone": "13800138000",
        "province": "广东省", "city": "深圳市",
        "district": "南山区", "detail": "科技园南区",
        "is_default": True
    })
    check("Create address", r.status_code == 201)

    # === 5. Get Categories ===
    r = session.get(f"{BASE_URL}/categories")
    check("Get categories", r.status_code == 200)

    # === 6. Get Products (empty) ===
    r = session.get(f"{BASE_URL}/products")
    check("Get products (empty)", r.status_code == 200)
    check("Empty list returned", len(r.json()["data"]["list"]) == 0)

    # === 7. Search Products (empty) ===
    r = session.get(f"{BASE_URL}/products/search?keyword=test")
    check("Search products", r.status_code == 200)

    # === 8. Get Banners ===
    r = session.get(f"{BASE_URL}/banners")
    check("Get banners", r.status_code == 200)

    # === 9. Get Payment Methods ===
    r = session.get(f"{BASE_URL}/payment/methods")
    check("Payment methods", r.status_code == 200)
    check("Has 3 methods", len(r.json()["data"]) == 3)

    # === 10. Get Active Flash Sales ===
    r = session.get(f"{BASE_URL}/flash-sales")
    check("Flash sales", r.status_code == 200)

    # === 11. Product CRUD (admin only) ===
    # Not tested here — requires admin role

    # === 12. Points ===
    r = session.get(f"{BASE_URL}/points")
    check("Query points", r.status_code == 200)
    check("Initial points = 0", r.json()["data"]["points"] == 0)

    r = session.get(f"{BASE_URL}/points/logs")
    check("Empty point logs", r.status_code == 200)

    # === 13. Favorites (empty) ===
    r = session.get(f"{BASE_URL}/favorites")
    check("Empty favorites", r.status_code == 200)

    # === 14. Browse History (empty) ===
    r = session.get(f"{BASE_URL}/browse-history")
    check("Empty browse history", r.status_code == 200)

    # === 15. Cart (empty) ===
    r = session.get(f"{BASE_URL}/cart")
    check("Empty cart", r.status_code == 200)
    r = session.get(f"{BASE_URL}/cart/count")
    check("Cart count 0", r.status_code == 200)

    # === 16. Orders (empty) ===
    r = session.get(f"{BASE_URL}/orders")
    check("Empty orders", r.status_code == 200)

    # === 17. Available Coupons ===
    r = session.get(f"{BASE_URL}/coupons/available?amount=10000")
    check("Available coupons", r.status_code == 200)

    # === 18. Logout / Unauthorized ===
    session.headers.pop("Authorization")
    r = session.get(f"{BASE_URL}/user/profile")
    check("Unauthorized access rejected", r.status_code == 401)

    print(f"\n{'='*50}")
    print(f"✅ All smoke tests passed!")
    print(f"{'='*50}")


if __name__ == "__main__":
    main()
```

- [ ] **Step 2: Create k6 load test script**

Create `server/scripts/k6_test.js`:

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export const options = {
  stages: [
    { duration: '10s', target: 5 },   // 逐步增加到 5 虚拟用户
    { duration: '30s', target: 5 },   // 保持 5 用户 30 秒
    { duration: '10s', target: 0 },   // 逐步减少
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% 请求在 2 秒内
    http_req_failed: ['rate<0.05'],    // 失败率低于 5%
  },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  // 1. 公开接口
  const healthRes = http.get('http://localhost:8080/health');
  check(healthRes, { 'health check ok': (r) => r.status === 200 });

  const bannerRes = http.get(`${BASE_URL}/banners`);
  check(bannerRes, { 'banners ok': (r) => r.status === 200 });

  const categoriesRes = http.get(`${BASE_URL}/categories`);
  check(categoriesRes, { 'categories ok': (r) => r.status === 200 });

  const productsRes = http.get(`${BASE_URL}/products`);
  check(productsRes, { 'products ok': (r) => r.status === 200 });

  const paymentMethodsRes = http.get(`${BASE_URL}/payment/methods`);
  check(paymentMethodsRes, { 'payment methods ok': (r) => r.status === 200 });

  // 2. 注册 + 登录
  const phone = `138${String(Date.now()).slice(-10)}`;
  const registerRes = http.post(`${BASE_URL}/user/register`, JSON.stringify({
    phone, password: 'test1234', nickname: 'k6_user',
  }), { headers: { 'Content-Type': 'application/json' } });
  check(registerRes, { 'register ok': (r) => r.status === 201 });

  const token = registerRes.json().data?.token;

  if (token) {
    const authHeaders = {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    };

    // 3. 认证接口
    const profileRes = http.get(`${BASE_URL}/user/profile`, authHeaders);
    check(profileRes, { 'profile ok': (r) => r.status === 200 });

    const pointsRes = http.get(`${BASE_URL}/points`, authHeaders);
    check(pointsRes, { 'points ok': (r) => r.status === 200 });

    const cartRes = http.get(`${BASE_URL}/cart`, authHeaders);
    check(cartRes, { 'cart ok': (r) => r.status === 200 });

    const ordersRes = http.get(`${BASE_URL}/orders`, authHeaders);
    check(ordersRes, { 'orders ok': (r) => r.status === 200 });

    const favoritesRes = http.get(`${BASE_URL}/favorites`, authHeaders);
    check(favoritesRes, { 'favorites ok': (r) => r.status === 200 });

    // 4. 创建地址
    const addrRes = http.post(`${BASE_URL}/addresses`, JSON.stringify({
      name: 'k6', phone: '13800138000',
      province: '广东', city: '深圳', district: '南山', detail: '测试地址',
      is_default: true,
    }), authHeaders);
    check(addrRes, { 'address created': (r) => r.status === 201 });
  }

  sleep(1);
}
```

- [ ] **Step 3: Start server with SQLite for Python smoke test**

Run: `cd server && go build -o server ./cmd/server && ./server &`
Expected: Server starts on :8080 (using dev config, SQLite)

- [ ] **Step 4: Run Python smoke tests**

Run: `pip install requests && python scripts/api_test.py`
Expected: All tests pass (✅ markers)

- [ ] **Step 5: Stop dev server**

Run: `kill %1`

- [ ] **Step 6: Start server with PostgreSQL for k6 load test**

`main.go` hardcodes `config/config.yaml`, but Viper already reads `GOSH_*` env vars. Use env overrides:

```bash
cd server
GOSH_DATABASE_DRIVER=postgres \
GOSH_DATABASE_HOST=127.0.0.1 \
GOSH_DATABASE_PORT=5432 \
GOSH_DATABASE_USER=gosh \
GOSH_DATABASE_PASSWORD=gosh \
GOSH_DATABASE_DBNAME=gosh \
./server &
```
Expected: Server starts on :8080 connected to PostgreSQL

- [ ] **Step 7: (Optional) Run k6 load test against PostgreSQL**

Run: `k6 run scripts/k6_test.js`
Expected: Tests complete, thresholds met

- [ ] **Step 8: Stop server and PostgreSQL**

Run: `kill %1 && docker compose down`
