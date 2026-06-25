# Phase 7 (7.3 + 7.4) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development or executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Implement order timeout configuration, Rebuy skipped-items response, cart max_buyable field, unified error code system, idempotent response completeness, and error code migration for core handlers.

**Architecture:** All changes are internal to the server package. Configuration is extended via viper. Error codes live in a new `pkg/errcode/` package. Handler error responses are gradually migrated from generic `-1` codes to typed error codes.

**Tech Stack:** Go 1.26, Gin, GORM, Viper, Zap

---

### Task 1: 订单超时配置化

**Files:**
- Modify: `server/config/config.yaml:41`
- Modify: `server/internal/config/config.go:17`
- Modify: `server/internal/scheduler/order_timeout.go:12-25`
- Modify: `server/cmd/server/main.go:85-86`

- [ ] **Step 1: Add `order` section to config.yaml**

In `server/config/config.yaml`, append at end (before EOF):
```yaml
order:
  timeout_minutes: 30
```

- [ ] **Step 2: Add `OrderConfig` struct to config.go**

In `server/internal/config/config.go`, after `LoggerConfig` struct, add:
```go
type OrderConfig struct {
	TimeoutMinutes int `mapstructure:"timeout_minutes"`
}
```

In `Config` struct, add field:
```go
Order OrderConfig `mapstructure:"order"`
```

- [ ] **Step 3: Update scheduler to accept timeout parameter**

In `server/internal/scheduler/order_timeout.go`, replace:
```go
const orderTimeout = 30 * time.Minute

type Scheduler struct {
	stopCh chan struct{}
}

func New() *Scheduler {
	return &Scheduler{
		stopCh: make(chan struct{}),
	}
}
```

With:
```go
type Scheduler struct {
	timeout time.Duration
	stopCh  chan struct{}
}

func New(timeoutMinutes int) *Scheduler {
	return &Scheduler{
		timeout: time.Duration(timeoutMinutes) * time.Minute,
		stopCh:  make(chan struct{}),
	}
}
```

In `Start` method, remove the hardcoded `zap.Duration("timeout", orderTimeout)` and replace with `zap.Duration("timeout", s.timeout)`.

In `cancelExpiredOrders`, replace `-orderTimeout` with `-s.timeout`. Change signature to `func (s *Scheduler) cancelExpiredOrders(log *zap.Logger)`.

- [ ] **Step 4: Update main.go to pass config**

In `server/cmd/server/main.go`, replace:
```go
orderScheduler := scheduler.New()
```

With:
```go
orderScheduler := scheduler.New(config.AppConfig.Order.TimeoutMinutes)
```

- [ ] **Step 5: Run tests**

Run: `cd server && go build ./...`
Expected: success

- [ ] **Step 6: Commit**

No commit — user will commit explicitly.

---

### Task 2: Rebuy 补充 — 返回 SkippedItems 而非 error

**Files:**
- Modify: `server/internal/dto/response/order.go:3-8`
- Modify: `server/internal/service/order/order.go:379-461`
- Modify: `server/internal/handler/order/order.go:314-335`

- [ ] **Step 1: Add RebuyResponse + SkippedItem types**

In `server/internal/dto/response/order.go`, after `OrderItemResponse` (line 43), add:
```go
type RebuyResponse struct {
	Cart         CartListResponse `json:"cart"`
	SkippedItems []SkippedItem    `json:"skipped_items"`
}

type SkippedItem struct {
	SKUID  uint   `json:"sku_id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type CartListResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCount int                `json:"total_count"`
}
```

- [ ] **Step 2: Modify Rebuy service to return structured response**

In `server/internal/service/order/order.go`, update `Rebuy` method to return `(*response.RebuyResponse, error)` instead of `error`:

```go
func (s *service) Rebuy(userID, orderID uint) (*response.RebuyResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotBelongToUser
	}
	if order.Status != model.OrderStatusCompleted {
		return nil, ErrInvalidOrderStatus
	}

	var items []model.OrderItem
	if err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
		return nil, err
	}

	skuIDs := make([]uint, len(items))
	for i, item := range items {
		skuIDs[i] = item.SKUID
	}

	var skus []model.ProductSKU
	if err := database.DB.Where("id IN ?", skuIDs).Find(&skus).Error; err != nil {
		return nil, err
	}
	skuMap := make(map[uint]model.ProductSKU, len(skus))
	for _, sku := range skus {
		skuMap[sku.ID] = sku
	}

	productIDs := make([]uint, 0)
	productIDSet := make(map[uint]bool)
	for _, sku := range skus {
		if !productIDSet[sku.ProductID] {
			productIDs = append(productIDs, sku.ProductID)
			productIDSet[sku.ProductID] = true
		}
	}

	var products []model.Product
	if err := database.DB.Select("id, name, status").Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	productStatus := make(map[uint]string, len(products))
	productNames := make(map[uint]string, len(products))
	for _, p := range products {
		productStatus[p.ID] = p.Status
		productNames[p.ID] = p.Name
	}

	var skipped []response.SkippedItem
	var cartItems []response.CartItemResponse
	for _, item := range items {
		sku, ok := skuMap[item.SKUID]
		if !ok {
			skipped = append(skipped, response.SkippedItem{
				SKUID:  item.SKUID,
				Name:   item.ProductName,
				Reason: "商品已删除",
			})
			continue
		}
		if productStatus[sku.ProductID] != model.ProductStatusOn {
			skipped = append(skipped, response.SkippedItem{
				SKUID:  item.SKUID,
				Name:   productNames[sku.ProductID],
				Reason: "商品已下架",
			})
			continue
		}

		var existingCart model.Cart
		err := database.DB.Where("user_id = ? AND sku_id = ?", userID, item.SKUID).First(&existingCart).Error
		if err == nil {
			database.DB.Model(&existingCart).Update("quantity", item.Quantity)
		} else {
			database.DB.Create(&model.Cart{
				UserID:   userID,
				SKUID:    item.SKUID,
				Quantity: item.Quantity,
				Selected: true,
			})
		}

		cartItems = append(cartItems, response.CartItemResponse{
			SKUID:    item.SKUID,
			Quantity: item.Quantity,
			Price:    sku.Price,
		})
	}

	return &response.RebuyResponse{
		Cart: response.CartListResponse{
			Items:      cartItems,
			TotalCount: len(cartItems),
		},
		SkippedItems: skipped,
	}, nil
}
```

Update the `Service` interface return type:
```go
Rebuy(userID, orderID uint) (*response.RebuyResponse, error)
```

- [ ] **Step 3: Update Rebuy handler**

In `server/internal/handler/order/order.go`, replace the Rebuy handler:
```go
func (h *Handler) Rebuy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	userID, _ := c.Get("user_id")
	resp, err := h.svc.Rebuy(userID.(uint), uint(id))
	if err != nil {
		if err == svc.ErrOrderNotFound || err == svc.ErrOrderNotBelongToUser {
			response.NotFound(c, err.Error())
			return
		}
		if err == svc.ErrInvalidOrderStatus {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "rebuy failed")
		return
	}
	response.Success(c, resp)
}
```

- [ ] **Step 4: Run tests**

Run: `cd server && go build ./... && go test ./internal/handler/order/ -v -count=1 2>&1 | tail -20`
Expected: all pass

---

### Task 3: Cart max_buyable 字段

**Files:**
- Modify: `server/internal/dto/response/cart.go:18`
- Modify: `server/internal/service/cart/cart.go:71-73, 330-362`

- [ ] **Step 1: Add MaxBuyable field to CartItemResponse**

In `server/internal/dto/response/cart.go`, after `Stock int` field (line 18), add:
```go
MaxBuyable int `json:"max_buyable"`
```

- [ ] **Step 2: Calculate MaxBuyable in toItemResponse**

In `server/internal/service/cart/cart.go`, add a constant at package level (after imports):
```go
const CartMaxQuantity = 99
```

In `toItemResponse` function, before the return statement, add:
```go
maxBuyable := sku.Stock
if maxBuyable > CartMaxQuantity {
	maxBuyable = CartMaxQuantity
}
```

Then add to the CartItemResponse literal:
```go
MaxBuyable: maxBuyable,
```

The field should be placed after `Stock` in the struct literal. The full return becomes:
```go
return &response.CartItemResponse{
	ID:          cart.ID,
	UserID:      cart.UserID,
	SKUID:       cart.SKUID,
	Quantity:    cart.Quantity,
	Selected:    cart.Selected,
	ProductName: productName,
	SKUName:     sku.Name,
	Image:       image,
	Price:       sku.Price,
	Stock:       sku.Stock,
	MaxBuyable:  maxBuyable,
	ProductID:   sku.ProductID,
	ProductOn:   online,
	CreatedAt:   cart.CreatedAt.Format("2006-01-02 15:04:05"),
}
```

- [ ] **Step 3: Run tests**

Run: `cd server && go build ./... && go test ./internal/handler/cart/ -v -count=1 2>&1 | tail -20`
Expected: all pass

---

### Task 4: 统一错误码基础包 + response 扩展

**Files:**
- Create: `server/pkg/errcode/errcode.go`
- Modify: `server/pkg/response/response.go:9-13`

- [ ] **Step 1: Create errcode package**

Create `server/pkg/errcode/errcode.go`:
```go
package errcode

// 通用 1xxx
const (
	ErrBadRequest   = 1001
	ErrUnauthorized = 1002
	ErrForbidden    = 1003
	ErrNotFound     = 1004
	ErrConflict     = 1005
	ErrInternal     = 1999
)

// 用户 2xxx
const (
	ErrUserExists    = 2001
	ErrPasswordWrong = 2002
	ErrUserNotFound  = 2003
)

// 商品 3xxx
const (
	ErrCategoryNotFound  = 3001
	ErrProductNotFound   = 3002
	ErrSKUNotFound       = 3003
	ErrInsufficientStock = 3004
)

// 订单 4xxx
const (
	ErrOrderNotFound   = 4001
	ErrOrderStatus     = 4002
	ErrCartEmpty       = 4003
)

// 支付 5xxx
const (
	ErrPaymentMethod = 5001
	ErrPaymentFailed = 5002
	ErrRefundFailed  = 5003
)

// 营销 6xxx
const (
	ErrCouponNotFound = 6001
	ErrCouponSoldOut  = 6002
	ErrCouponReceived = 6003
)
```

- [ ] **Step 2: Add ErrorWithCode to response package**

In `server/pkg/response/response.go`, after `Error` function, add:
```go
func ErrorWithCode(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequestWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusBadRequest, code, message)
}

func NotFoundWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusNotFound, code, message)
}

func UnauthorizedWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusUnauthorized, code, message)
}

func ForbiddenWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusForbidden, code, message)
}

func InternalErrorWithCode(c *gin.Context, code int, message string) {
	ErrorWithCode(c, http.StatusInternalServerError, code, message)
}
```

- [ ] **Step 3: Run tests**

Run: `cd server && go build ./...`
Expected: success, no breakage (old functions still work)

---

### Task 5: 统一错误码迁移核心 handler

**Files:**
- Modify: `server/internal/handler/order/order.go`
- Modify: `server/internal/handler/cart/cart.go`
- Modify: `server/internal/handler/user/user.go`
- Modify: `server/internal/handler/product/product.go`
- Modify: `server/internal/handler/payment/payment.go`

- [ ] **Step 1: Migrate order handler**

In `server/internal/handler/order/order.go`, add import:
```go
"gosh/pkg/errcode"
```

Update each handler function:

**Create:** Replace `response.BadRequest(c, err.Error())` with appropriate code. For cart empty → `errcode.ErrCartEmpty`, insufficient stock → `errcode.ErrInsufficientStock`. Example:
```go
if errors.Is(err, svc.ErrCartEmpty) {
    response.BadRequestWithCode(c, errcode.ErrCartEmpty, err.Error())
    return
}
if errors.Is(err, svc.ErrInsufficientStock) {
    response.BadRequestWithCode(c, errcode.ErrInsufficientStock, err.Error())
    return
}
```

**GetByID:** Replace `response.NotFound(c, err.Error())` → `response.NotFoundWithCode(c, errcode.ErrOrderNotFound, err.Error())`

**Cancel:** Replace `response.NotFound` → `NotFoundWithCode(c, errcode.ErrOrderNotFound, ...)`, `response.BadRequest` → `BadRequestWithCode(c, errcode.ErrOrderStatus, ...)`

**Pay:** Same pattern as Cancel.

**Ship:** Same pattern as Cancel.

**Confirm:** Same pattern as Cancel.

**ApplyPoints:** Replace `response.BadRequest` (for insufficient points) with `BadRequestWithCode(c, errcode.ErrBadRequest, ...)`.

**Rebuy:** Replace `response.NotFound` → `NotFoundWithCode(c, errcode.ErrOrderNotFound, ...)`, `response.BadRequest` → `BadRequestWithCode(c, errcode.ErrOrderStatus, ...)`

- [ ] **Step 2: Migrate cart handler**

In `server/internal/handler/cart/cart.go`, add imports:
```go
"errors"
"gosh/pkg/errcode"
"gosh/internal/service/cart"
```

Replace:
- `response.BadRequest(c, err.Error())` → `response.BadRequestWithCode(c, errcode.ErrBadRequest, err.Error())` (for binding errors)
- `response.NotFound(c, err.Error())` → `response.NotFoundWithCode(c, errcode.ErrNotFound, err.Error())` (cart not found)
- SKU not found → `response.BadRequestWithCode(c, errcode.ErrSKUNotFound, err.Error())`
- Product off shelf → `response.BadRequestWithCode(c, errcode.ErrProductNotFound, err.Error())`

- [ ] **Step 3: Migrate user handler**

In `server/internal/handler/user/user.go`, add import `"gosh/pkg/errcode"`.

Replace error responses:
- Login password wrong → `response.BadRequestWithCode(c, errcode.ErrPasswordWrong, ...)`
- User already exists → `response.ErrorWithCode(c, 409, errcode.ErrUserExists, ...)`
- User not found → `response.NotFoundWithCode(c, errcode.ErrUserNotFound, ...)`
- Unauthorized → `response.UnauthorizedWithCode(c, errcode.ErrUnauthorized, ...)`
- Forbidden → `response.ForbiddenWithCode(c, errcode.ErrForbidden, ...)`

- [ ] **Step 4: Migrate product handler**

In `server/internal/handler/product/product.go`, add import `"gosh/pkg/errcode"`.

Replace error responses:
- Category not found → `response.NotFoundWithCode(c, errcode.ErrCategoryNotFound, ...)`
- Product not found → `response.NotFoundWithCode(c, errcode.ErrProductNotFound, ...)`

- [ ] **Step 5: Migrate payment handler**

In `server/internal/handler/payment/payment.go`, add import `"gosh/pkg/errcode"`.

Replace error responses:
- Payment method not supported → `BadRequestWithCode(c, errcode.ErrPaymentMethod, ...)`
- Payment failed → `InternalErrorWithCode(c, errcode.ErrPaymentFailed, ...)`
- Refund failed → `InternalErrorWithCode(c, errcode.ErrRefundFailed, ...)`

- [ ] **Step 6: Run tests**

Run: `cd server && go test ./... 2>&1 | grep -E "^(ok|FAIL|---|\?)" `
Expected: all pass (old response functions still exist, so no breakage)

---

### Task 6: 幂等性响应返回完整 data

**Files:**
- Modify: `server/internal/service/order/order.go:64-69`

- [ ] **Step 1: Store order response JSON in IdempotencyRecord on creation**

In `server/internal/service/order/order.go`, after the order is created successfully (after line 196 `resp := response.ToOrderResponse(fullOrder)`), store the response:
```go
// Store response for idempotency
respJSON, _ := json.Marshal(resp)
s.orderRepo.CreateIdempotency(&model.IdempotencyRecord{
    Key:      idempotentKey,
    Response: string(respJSON),
})
```

Remove the existing idempotency record creation at line 89:
```go
// DELETE this line:
if err := tx.Create(&model.IdempotencyRecord{Key: idempotentKey}).Error; err != nil {
```

Add `"encoding/json"` to imports.

- [ ] **Step 2: Return stored response on idempotent hit**

Replace the idempotent check block (lines 65-69):
```go
// 幂等检查
existing, err := s.orderRepo.FindIdempotency(idempotentKey)
if err == nil && existing != nil {
    var resp response.OrderResponse
    if existing.Response != "" {
        json.Unmarshal([]byte(existing.Response), &resp)
    }
    return &resp, nil
}
```

With:
```go
// 幂等检查
existing, err := s.orderRepo.FindIdempotency(idempotentKey)
if err == nil && existing != nil && existing.Response != "" {
    var resp response.OrderResponse
    if err := json.Unmarshal([]byte(existing.Response), &resp); err == nil {
        return &resp, nil
    }
}
```

Note: still need the `var resp response.OrderResponse` for the fallback after the block.

- [ ] **Step 3: Run tests**

Run: `cd server && go test ./internal/handler/order/ -v -count=1 2>&1 | tail -20`
Expected: all pass

- [ ] **Step 4: Build check**

Run: `cd server && go build ./...`
Expected: success

---

## Self-Review

- **Spec coverage:** All items from the spec are covered: 7.3.1 (Task 1), 7.3.3 (Task 2), 7.4.3 (Task 3), 7.4.1 (Task 4+5), 7.4.2 (Task 6). 7.3.2 (inventory strategy) was decided to keep as-is.
- **Placeholder scan:** No TBD/TODO/fill-in-later patterns. All code is concrete.
- **Type consistency:** `RebuyResponse` is defined in Task 2 and used consistently. `CartMaxQuantity` constant defined in Task 3 and used in same task. Error codes in Task 4 match usage in Task 5.
- **Task ordering:** Tasks are independent (1, 2, 3, 4+5 depend on 4) and ordered correctly.
