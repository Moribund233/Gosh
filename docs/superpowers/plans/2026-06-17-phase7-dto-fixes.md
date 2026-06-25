# Phase 7 — DTO 修正 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Apply 4 DTO fixes from the Phase 7 spec: Images `[]string`, SelectRequest json tag, CreateOrderRequest optional items, CalculateCouponRequest auto-best.

**Architecture:** Each fix is self-contained in DTO + service layers. GORM's `serializer:json` handles `[]string` persistence. The coupon auto-best uses the existing `GetAvailable` path internally.

**Tech Stack:** Go 1.26, Gin, GORM, SQLite (test)

**Run all tests:** `cd server && go test -count=1 ./internal/...`

---
### Task A: Images string → []string

Change `Images` from `string` to `[]string` across model, DTOs, and service code that reads `product.Images` for order/cart snapshots.

**Files to Modify:**
- `internal/model/product.go:13,45` – `Images string` → `Images []string` with `gorm:"type:text;serializer:json"`
- `internal/dto/request/product.go:11,35,59` – `Images string` → `Images []string`
- `internal/dto/response/product.go:17,43` – `Images string` → `Images []string`
- `internal/service/product/product.go:113-114` – `req.Images != ""` → `len(req.Images) > 0`
- `internal/service/order/order.go:136` – `product.Images` (now `[]string`) → pick first image
- `internal/service/cart/cart.go:337,343` – `product.Images` / `p.Images` → pick first image

- [ ] **Step 1: Update model**

In `internal/model/product.go`:

```go
type Product struct {
    ...
    Images  []string `gorm:"type:text;serializer:json" json:"images"`
    ...
}

type ProductReview struct {
    ...
    Images  []string `gorm:"type:text;serializer:json" json:"images"`
    ...
}
```

- [ ] **Step 2: Update request DTOs**

In `internal/dto/request/product.go`:

```go
Images []string `json:"images"`
```

On `CreateProductRequest` (line 11), `UpdateProductRequest` (line 35), `CreateReviewRequest` (line 59).

- [ ] **Step 3: Update response DTOs**

In `internal/dto/response/product.go`:

```go
Images []string `json:"images"`
```

On `ProductResponse` (line 17), `ReviewResponse` (line 43).

- [ ] **Step 4: Fix product service condition**

In `internal/service/product/product.go`, line 113:
```go
// Change from:
if req.Images != "" {
    product.Images = req.Images
}
// To:
if len(req.Images) > 0 {
    product.Images = req.Images
}
```

- [ ] **Step 5: Fix order service image snapshot**

In `internal/service/order/order.go`, replace line 136:

```go
// Change from:
Image: product.Images,
// To:
Image: firstImage(product.Images),
```

At the top of the file or as a package-level helper:

```go
func firstImage(images []string) string {
    if len(images) > 0 {
        return images[0]
    }
    return ""
}
```

- [ ] **Step 6: Fix cart service image snapshot**

In `internal/service/cart/cart.go`, lines 337 and 343:

```go
// Both change from:
image = product.Images   // line 337
image = p.Images         // line 343
// To:
image = firstImage(product.Images)
image = firstImage(p.Images)
```

Add the same `firstImage` helper (or share it — but standalone package func is simplest in each file).

- [ ] **Step 7: Run tests**

```bash
cd /home/alucard/文档/Project/go/gosh/server && go test -count=1 ./internal/...
```

Expected: all existing tests pass (images were already stored as serialized JSON in tests? No — they were stored as plain strings. Tests may have used `Images: "url1,url2"` as a string, which will no longer compile. Need to update test fixtures.)

**Important:** Find and update test fixtures that assign `Images: "some_string"` to use `Images: []string{"url1", "url2"}`.

Run grep to find all assignments to `Images` in test files:
```bash
cd /home/alucard/文档/Project/go/gosh/server && rg 'Images:\s*"' --include='*_test.go' -l
```

If any found, update them to `Images: []string{...}`.

---

### Task B: SelectRequest json tag — `"select"` → `"selected"`

**Files to Modify:**
- `internal/dto/request/cart.go:16` – json tag change
- `internal/service/cart/cart.go:267,274` – `req.Select` stays as Go field name, no change needed

- [ ] **Step 1: Change json tag**

In `internal/dto/request/cart.go`, line 16:

```go
// From:
Select bool `json:"select"`
// To:
Select bool `json:"selected"`
```

- [ ] **Step 2: Verify service code**

The service at `internal/service/cart/cart.go:267,274` uses `req.Select` as the Go field name. Since only the json tag changes, no Go code changes.

- [ ] **Step 3: Check test files for JSON key `"select"`**

```bash
cd /home/alucard/文档/Project/go/gosh/server && rg --include='*_test.go' '"select"' -l
```

If any test JSON payloads use `"select": true`, update to `"selected": true`.

- [ ] **Step 4: Run tests**

```bash
cd /home/alucard/文档/Project/go/gosh/server && go test -count=1 ./internal/...
```

---

### Task C: CreateOrderRequest — optional Items field

**Files to Modify:**
- Add: `internal/dto/request/order.go` — `CreateOrderItem` struct + `Items` field on `CreateOrderRequest`
- Modify: `internal/service/order/order.go` — in `Create`, read cart from `req.Items` when present, else from DB

- [ ] **Step 1: Add CreateOrderItem and Items field**

In `internal/dto/request/order.go`, add before `CreateOrderRequest`:

```go
type CreateOrderItem struct {
    SKUID    uint `json:"sku_id" binding:"required"`
    Quantity int  `json:"quantity" binding:"required,min=1"`
}
```

Then add `Items` field to `CreateOrderRequest`:

```go
type CreateOrderRequest struct {
    AddressID      uint              `json:"address_id"`
    Remark         string            `json:"remark" binding:"omitempty,max=200"`
    DeliveryMethod string            `json:"delivery_method" binding:"omitempty,oneof=standard express"`
    Items          []CreateOrderItem `json:"items,omitempty"`
}
```

- [ ] **Step 2: Modify order service Create logic**

In `internal/service/order/order.go`, replace the cart-fetching block (lines 76-82):

**From:**
```go
carts, err := s.cartRepo.FindSelectedByUserID(userID)
if err != nil {
    return nil, err
}
if len(carts) == 0 {
    return nil, ErrCartEmpty
}
```

**To:**
```go
// Use explicitly provided items, or fall back to selected cart items
var carts []model.Cart
if len(req.Items) > 0 {
    for _, item := range req.Items {
        carts = append(carts, model.Cart{
            SKUID:    item.SKUID,
            Quantity: item.Quantity,
            Selected: true,
        })
    }
} else {
    carts, err = s.cartRepo.FindSelectedByUserID(userID)
    if err != nil {
        return nil, err
    }
    if len(carts) == 0 {
        return nil, ErrCartEmpty
    }
}
```

- [ ] **Step 3: Fix the "clear cart" branch**

At line 172-174, we clear selected cart items. This should only happen when items came from the cart (not from `req.Items`). Wrap with a condition:

**From:**
```go
// 清空已购购物车
if err := tx.Where("user_id = ? AND selected = ?", userID, true).Delete(&model.Cart{}).Error; err != nil {
    return err
}
```

**To:**
```go
// 清空已购购物车 (only when items came from cart)
if len(req.Items) == 0 {
    if err := tx.Where("user_id = ? AND selected = ?", userID, true).Delete(&model.Cart{}).Error; err != nil {
        return err
    }
}
```

- [ ] **Step 4: Run tests**

```bash
cd /home/alucard/文档/Project/go/gosh/server && go test -count=1 ./internal/...
```

---

### Task D: CalculateCouponRequest — optional CouponID (auto-best)

**Files to Modify:**
- `internal/dto/request/coupon.go:16` — `CouponID` from `uint` to `*uint`, remove `required`
- `internal/dto/response/coupon.go:34-38` — Add `CouponID` field to response
- `internal/service/coupon/coupon.go:27` — `Calculate` interface: add `userID` param
- `internal/service/coupon/coupon.go:152` — Update signature + auto-best logic
- `internal/handler/coupon/coupon.go:117-136` — Pass `userID` to Calculate

- [ ] **Step 1: Update request DTO**

In `internal/dto/request/coupon.go`, line 15-16:

```go
type CalculateCouponRequest struct {
    OrderAmount int64 `json:"order_amount" binding:"required,min=0"`
    CouponID    *uint `json:"coupon_id,omitempty"`
}
```

- [ ] **Step 2: Update response DTO**

In `internal/dto/response/coupon.go`, add `CouponID` to `CouponCalculateResponse`:

```go
type CouponCalculateResponse struct {
    CouponID       *uint `json:"coupon_id,omitempty"`
    OriginalAmount int64 `json:"original_amount"`
    DiscountAmount int64 `json:"discount_amount"`
    PayAmount      int64 `json:"pay_amount"`
}
```

- [ ] **Step 3: Update service interface + implementation**

In `internal/service/coupon/coupon.go`, line 27:
```go
// From:
Calculate(req *request.CalculateCouponRequest) (*response.CouponCalculateResponse, error)
// To:
Calculate(userID uint, req *request.CalculateCouponRequest) (*response.CouponCalculateResponse, error)
```

Replace the entire `Calculate` method body (lines 152-191):

```go
func (s *service) Calculate(userID uint, req *request.CalculateCouponRequest) (*response.CouponCalculateResponse, error) {
    if req.CouponID != nil {
        // Specific coupon: calculate as before
        coupon, err := s.repo.FindByID(*req.CouponID)
        if err != nil {
            return nil, ErrCouponNotFound
        }

        now := time.Now()
        if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
            return nil, ErrCouponExpired
        }
        if coupon.Status != model.CouponStatusActive {
            return nil, ErrCouponExpired
        }
        if req.OrderAmount < coupon.Condition {
            return nil, ErrCouponNotApplicable
        }

        discount := calculateDiscount(req.OrderAmount, coupon)
        return &response.CouponCalculateResponse{
            CouponID:       &coupon.ID,
            OriginalAmount: req.OrderAmount,
            DiscountAmount: discount,
            PayAmount:      req.OrderAmount - discount,
        }, nil
    }

    // Auto-best: find the best available coupon
    available, err := s.GetAvailable(userID, req.OrderAmount)
    if err != nil {
        return nil, err
    }

    bestCouponID := (*uint)(nil)
    bestDiscount := int64(0)

    // Also fetch full coupon objects for calculation
    coupons, err := s.repo.FindActive(req.OrderAmount)
    if err != nil {
        return nil, err
    }
    couponMap := make(map[uint]*model.Coupon)
    for i := range coupons {
        couponMap[coupons[i].ID] = &coupons[i]
    }

    for _, uc := range available {
        c, ok := couponMap[uc.CouponID]
        if !ok {
            continue
        }
        if req.OrderAmount < c.Condition {
            continue
        }
        d := calculateDiscount(req.OrderAmount, c)
        if d > bestDiscount {
            bestDiscount = d
            id := c.ID
            bestCouponID = &id
        }
    }

    return &response.CouponCalculateResponse{
        CouponID:       bestCouponID,
        OriginalAmount: req.OrderAmount,
        DiscountAmount: bestDiscount,
        PayAmount:      req.OrderAmount - bestDiscount,
    }, nil
}

func calculateDiscount(amount int64, coupon *model.Coupon) int64 {
    switch coupon.Type {
    case model.CouponTypeFullReduce:
        if coupon.Discount > amount {
            return amount
        }
        return coupon.Discount
    case model.CouponTypeDiscount:
        d := amount * (100 - coupon.Discount) / 100
        if d > amount {
            return amount
        }
        return d
    }
    return 0
}
```

- [ ] **Step 4: Update handler**

In `internal/handler/coupon/coupon.go`, update `Calculate` method to pass `userID`:

```go
func (h *Handler) Calculate(c *gin.Context) {
    var req request.CalculateCouponRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    userID, _ := c.Get("user_id")
    resp, err := h.svc.Calculate(userID.(uint), &req)
    // ... rest stays the same
}
```

- [ ] **Step 5: Run tests**

```bash
cd /home/alucard/文档/Project/go/gosh/server && go test -count=1 ./internal/...
```

---

## Self-Review Checklist

- [x] **Spec coverage:** All 4 DTO fixes from Phase 7 spec covered (Images `[]string`, SelectRequest json tag, CreateOrderRequest items, CalculateCouponRequest auto-best)
- [x] **No placeholders:** All code blocks are complete, no TBD/TODO patterns
- [x] **Type consistency:** `[]string` flows through model→dto→service consistently. `firstImage` helper signature matches usage. `*uint` for optional CouponID is used consistently.
