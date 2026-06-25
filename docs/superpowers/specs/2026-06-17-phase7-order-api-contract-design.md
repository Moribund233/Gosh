# 阶段七：订单优化与 API 契约一致性

## 范围

仅覆盖 roadmap 阶段七中 **7.3（订单与库存策略优化）** 和 **7.4（API 契约一致性）** 两个子项。

---

## 7.3 订单与库存策略优化

### 7.3.1 订单超时配置化

**现状：** `scheduler/order_timeout.go` 中 `const orderTimeout = 30 * time.Minute` 硬编码。

**目标：** 可配置。

**改动：**

1. `config/config.yaml` 新增：
   ```yaml
   order:
     timeout_minutes: 30
   ```

2. `config/config.go` 新增字段：
   ```go
   type Config struct {
       // ...
       Order struct {
           TimeoutMinutes int `mapstructure:"timeout_minutes"`
       } `mapstructure:"order"`
   }
   ```

3. `scheduler/order_timeout.go` 移除硬编码常量，改为接收参数或读取 `config.AppConfig.Order.TimeoutMinutes`：
   ```go
   // Before
   const orderTimeout = 30 * time.Minute

   // After
   timeout := time.Duration(config.AppConfig.Order.TimeoutMinutes) * time.Minute
   ```

4. 依赖注入：scheduler 初始化时传入 timeout，或在定时器循环内每次从 config 读取。

### 7.3.2 库存扣减策略

**决策：** 保持现状——下单即扣（乐观锁 `UPDATE ... WHERE stock >= ? AND version = ?`），超时回滚。不引入 reserved_stock / actual_stock 两阶段模型。

### 7.3.3 Rebuy 补充

**现状：** `service/order/order.go:Rebuy` 可能因已下架/已删除 SKU 报 500。

**目标：** 优雅降级——跳过无效 SKU，在响应中返回 `skipped_items` 提示。

**改动：**

1. Rebuy 流程中检查 `product.status != 'on'` 或 SKU 已删除：
   - 跳过该商品
   - 收集到 `skipped []SkippedItem` 列表
   - 继续处理剩余有效商品

2. 定义响应结构：
   ```go
   type RebuyResponse struct {
       Cart          CartListResponse `json:"cart"`
       SkippedItems  []SkippedItem    `json:"skipped_items"`
   }

   type SkippedItem struct {
       SKUID    uint   `json:"sku_id"`
       Name     string `json:"name"`
       Reason   string `json:"reason"` // "已下架" 或 "已删除"
   }
   ```

---

## 7.4 API 契约一致性

### 7.4.1 统一错误码

**现状：** `pkg/response/response.go` 仅 `Code: 0`（成功）和 `Code: -1`（错误），客户端无法区分业务错误类型。

**目标：** 按模块分段的枚举错误码。

**设计：**

新建 `pkg/errcode/` 包：

```go
package errcode

// 通用 1xxx
const (
    ErrBadRequest     = 1001  // 参数错误
    ErrUnauthorized   = 1002  // 未登录/Token 过期
    ErrForbidden      = 1003  // 无权限
    ErrNotFound       = 1004  // 资源不存在
    ErrConflict       = 1005  // 冲突（如重复操作）
    ErrInternal       = 1999  // 系统内部错误
)

// 用户 2xxx
const (
    ErrUserExists     = 2001  // 手机号已注册
    ErrPasswordWrong  = 2002  // 密码错误
    ErrUserNotFound   = 2003  // 用户不存在
)

// 商品 3xxx
const (
    ErrCategoryNotFound = 3001
    ErrProductNotFound  = 3002
    ErrSKUNotFound      = 3003
    ErrInsufficientStock = 3004
)

// 订单 4xxx
const (
    ErrOrderNotFound   = 4001
    ErrOrderStatus     = 4002  // 状态不允许操作
    ErrCartEmpty       = 4003
)

// 支付 5xxx
const (
    ErrPaymentMethod   = 5001
    ErrPaymentFailed   = 5002
    ErrRefundFailed    = 5003
)

// 营销 6xxx
const (
    ErrCouponNotFound  = 6001
    ErrCouponSoldOut   = 6002
    ErrCouponReceived  = 6003  // 已领取
)
```

`pkg/response/` 新增 `ErrorWithCode`：

```go
func ErrorWithCode(c *gin.Context, httpStatus int, code int, msg string) {
    c.AbortWithStatusJSON(httpStatus, Response{
        Code:    code,
        Message: msg,
    })
}
```

**迁移策略：** 渐进式替换。
- 先定义 `errcode` 包 + `ErrorWithCode`
- 逐个 handler/service 替换 `response.Error` → `response.ErrorWithCode`（或保留别名 `response.ErrBadRequest(c, msg)`）
- 不要求在本次全部替换完成，但核心流程（订单/支付/购物车/用户）必须覆盖

### 7.4.2 幂等性响应

**现状：** `service/order/order.go` 幂等拦截命中后返回空 `data`。

**目标：** 返回与首次创建一致的完整 order response。

**改动：**

1. 幂等检查命中时，调用 `GetByID` 获取已有订单完整数据
2. 返回格式与首次创建一致（order detail + order items）
3. HTTP Status 保持 201（与首次一致）或 200（幂等语义）

### 7.4.3 购物车库存实时性

**现状：** `CartItemResponse` 返回 `stock`（SKU 总库存），但前端需要指导用户最大可买数量。

**目标：** 新增 `max_buyable` 字段。

**改动：**

1. `dto/response/cart.go` 新增字段：
   ```go
   MaxBuyable int `json:"max_buyable"`
   ```

2. `service/cart/cart.go` 组装时计算：
   ```go
   maxBuyable := sku.Stock
   if maxBuyable > CartMaxQuantity {
       maxBuyable = CartMaxQuantity // 购物车上限 99
   }
   ```

3. 同时定义常量 `CartMaxQuantity = 99` 供 7.2 购物车数量校验复用

---

## 不在此次范围

- 7.1 DTO 修正（Images / select / items / calculate_coupon）
- 7.2 安全增强（购物车上限、上传限制、XSS）
- 7.5 可观测性（request_id / 慢查询 / 结构化日志）
- 7.6 CI/CD

以上项将在阶段七后段按 roadmap 顺序推进。

---

## 文件改动清单

| 文件 | 改动 |
|------|------|
| `server/config/config.yaml` | 新增 `order.timeout_minutes` |
| `server/internal/config/config.go` | 新增 `Order.TimeoutMinutes` |
| `server/internal/scheduler/order_timeout.go` | 硬编码 → 配置化 |
| `server/internal/service/order/order.go` | Rebuy 跳过提示 + 幂等返回完整 data |
| `server/pkg/errcode/errcode.go` | 新建，错误码常量 |
| `server/pkg/response/response.go` | 新增 `ErrorWithCode` |
| `server/internal/dto/response/cart.go` | 新增 `MaxBuyable` |
| `server/internal/service/cart/cart.go` | 计算 `max_buyable` |
| `server/internal/handler/*.go` | 渐进式替换 `Error` → `ErrorWithCode` |
