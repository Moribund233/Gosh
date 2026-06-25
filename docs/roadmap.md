# Gosh Mall 开发路线图

> 开发原则：**TDD** — 所有功能开发严格遵循 红→绿→重构 循环。
> 每个功能单元：先写测试 → 再写实现 → 最后重构。
>
> **开发策略**：后端先行，API 文档后出，多客户端并行。
> 阶段一～五 全栈开发 → 阶段六 API 文档（Swagger）→ 阶段七 多客户端前端。
> 后端接口设计保持 RESTful 风格统一，待 Swagger 就绪后各客户端可并行开发。

## 阶段一：基础搭建（第 1-2 周）

- [x] Go 模块初始化
- [x] uniapp 前端项目初始化
- [x] 后端目录结构设计 (`server/`)
- [x] uniapp 原型设计
- [x] 后端核心框架搭建
  - [x] HTTP 框架（Gin）
  - [x] 配置管理（Viper + YAML）
  - [x] 数据库连接（GORM + PostgreSQL / SQLite）
  - [x] 日志系统（Zap）
  - [x] 统一响应格式
- [x] 前端基础框架
  - [x] TabBar 与页面路由
  - [x] TypeScript 集成
  - [x] uView Plus UI 库
  - [x] HTTP 请求封装
  - [x] Pinia 状态管理
  - [x] 设计令牌映射到 uni.scss
  - [ ] ~~原型标注与页面联调~~ → **移至阶段七多客户端前端**，待 API 定型后统一进行（见阶段七说明）
- [x] 测试基础建设
  - [x] 测试框架（testify + go-sqlmock）
  - [x] 测试辅助函数（DB mock、HTTP mock）
  - [x] 测试目录结构规范

## 阶段二：用户系统（TDD）

### 角色定义

| 角色 | 标识 | 权限范围 |
|------|------|---------|
| **超级管理员** | `super_admin` | 全权限：角色管理、平台配置、所有数据 |
| **运营** | `operator` | 商品管理、订单管理、用户管理 |
| **商户** | `merchant` | 管理自有店铺：商品、订单、数据 |
| **客服** | `support` | 订单查看、售后处理、用户咨询 |
| **普通用户** | `user` | 前台注册、浏览、购物、下单 |

数据隔离：所有业务表加 `tenant_id` / `shop_id`，商户仅可见自身数据。

### 第 3 周 — 用户认证 ✅

- [x] **测试先行**：用户注册接口测试
  - [x] 正常注册返回 token
  - [x] 手机号已注册返回 409
  - [x] 参数缺失返回 400
- [x] 用户模型 + 密码加密（bcrypt）
- [x] 用户注册接口实现（默认 `user` 角色）
- [x] **测试先行**：登录接口测试
  - [x] 正确凭证返回 token
  - [x] 错误密码返回 401
- [x] JWT 签发（payload 含 `role` + `tenant_id`）
- [x] 登录接口实现
- [x] **测试先行**：JWT 鉴权中间件测试
  - [x] 有效 token 通过
  - [x] 过期 token 返回 401
  - [x] 无 token 返回 401
- [x] **测试先行**：角色校验中间件测试
  - [x] 管理员接口拒绝普通用户
- [x] 角色权限中间件实现

### 第 4 周 — 用户中心与权限 ✅

- [x] **测试先行**：用户个人信息接口测试
- [x] 个人中心接口实现（更新昵称/头像）
- [x] **测试先行**：收货地址 CRUD 测试
- [x] 收货地址管理实现（Create/List/Update/Delete）
- [x] **测试先行**：收藏夹接口测试
- [x] 收藏夹（增删改查）实现
- [x] **测试先行**：浏览记录接口测试
- [x] 浏览记录实现
- [x] **测试先行**：角色 CRUD 测试
- [x] 角色管理实现（超级管理员维护）
- [x] **测试先行**：商户入驻申请测试
  - [x] 提交申请
  - [x] 管理员审核通过/驳回
  - [x] 审核后角色变更为 `merchant`
- [x] 商户入驻申请与审核流程

## 阶段三：商品系统（TDD）✅

### 第 5 周 — 商品模型与分类 ✅

- [x] **测试先行**：商品分类测试（多级树形结构）
- [x] 分类管理实现（CRUD + 树形，含分类横幅）
- [x] **测试先行**：SPU / SKU 模型测试
- [x] SPU / SKU 模型与数据库迁移
- [x] **测试先行**：商品列表 API 测试
  - [x] 分页返回
  - [x] 按分类筛选
  - [x] 按标签筛选（热卖、新品、精选等）
  - [x] 空分类返回空列表
- [x] 商品列表接口实现（含标签过滤、排序）

### 第 6 周 — 商品功能与详情 ✅

- [x] **测试先行**：商品详情 API 测试
  - [x] 基础信息（名称、价格、销量、标签）
  - [x] 图片列表 / 轮播
  - [x] 规格参数（产地、保质期等）
  - [x] SKU 多规格
- [x] 商品详情接口实现
- [x] **测试先行**：商品评价接口测试
- [x] 商品评价实现（评分、内容、脱敏用户名）
- [x] **测试先行**：商品搜索测试（模糊匹配）
- [x] 商品搜索实现（LIKE 模糊匹配）
  - [x] 热搜排行
  - [x] 搜索历史记录（用户级）
- [x] **测试先行**：图片上传测试
- [x] 图片/素材上传实现（multipart + base64，URL 返回）
- [x] **测试先行**：轮播/品牌故事管理测试
- [x] 首页轮播横幅 CRUD + 品牌故事管理实现
- [x] **测试先行**：后台商品管理 CRUD 测试
- [x] 后台商品管理实现

## 阶段四：购物车与订单（TDD）

> ⚠️ **核心资金流安全原则**
> - 所有价格从服务端 DB 读取，**绝不信任客户端**传入的 price/subtotal/total
> - 库存扣减使用**乐观锁**（`UPDATE ... WHERE stock >= ? AND version = ?`），防超卖
> - 订单创建在**数据库事务**中执行，任何一步失败 → ROLLBACK
> - POST /orders 强制要求 **Idempotent-Key** 请求头，防重复下单
> - 订单取消/过期时**补偿恢复库存**
> - 所有订单状态变更写入 **OrderLog 审计表**

### 数据模型

```go
// Cart — 购物车（唯一约束: user_id + sku_id）
type Cart struct {
    BaseModel
    UserID    uint `gorm:"uniqueIndex:idx_cart_user_sku;not null"`
    SKUID     uint `gorm:"uniqueIndex:idx_cart_user_sku;not null"`
    Quantity  int  `gorm:"not null;default:1"`
    Selected  bool `gorm:"default:true"`
}

// Order — 订单（order_no 唯一索引）
type Order struct {
    BaseModel
    OrderNo        string     `gorm:"uniqueIndex;size:32;not null"`   // YYYYMMDDHHmmss + 4随机 + 2校验 = 20位
    UserID         uint       `gorm:"index;not null"`
    Status         string     `gorm:"size:20;default:pending_payment;index"`
    TotalAmount    int64      `gorm:"not null"`                       // 商品总价（分）
    ShippingFee    int64      `gorm:"default:0"`                      // 运费
    DiscountAmount int64      `gorm:"default:0"`                      // 优惠减免
    PayAmount      int64      `gorm:"not null"`                       // 实付 = TotalAmount + ShippingFee - DiscountAmount
    Remark         string     `gorm:"size:200"`
    DeliveryMethod string     `gorm:"size:32"`                        // standard, express
    // 收货地址快照（下单时复制，非 FK 引用）
    AddressName    string     `gorm:"size:32"`
    AddressPhone   string     `gorm:"size:20"`
    AddressDetail  string     `gorm:"size:512"`                       // JSON 格式
    // 各状态时间戳
    PaidAt         *time.Time
    ShippedAt      *time.Time
    CompletedAt    *time.Time
    CancelledAt    *time.Time
    CancelReason   string     `gorm:"size:128"`
    Version        int        `gorm:"default:0"`                      // 乐观锁
}

// OrderItem — 订单明细（价格/名称/图片均为下单时快照）
type OrderItem struct {
    BaseModel
    OrderID     uint   `gorm:"index;not null"`
    SKUID       uint   `gorm:"not null"`
    ProductName string `gorm:"size:128;not null"`                    // 快照
    SKUName     string `gorm:"size:64"`                              // 快照
    Image       string `gorm:"size:256"`                             // 快照
    Price       int64  `gorm:"not null"`                             // 快照（分），服务端读取，绝不信任客户端
    Quantity    int    `gorm:"not null"`
    Subtotal    int64  `gorm:"not null"`                             // Price × Quantity，服务端计算
}

// Payment — 支付记录（按事务幂等设计）
type Payment struct {
    BaseModel
    OrderID       uint       `gorm:"index;not null"`
    OrderNo       string     `gorm:"size:32;not null"`
    Method        string     `gorm:"size:20;not null"`              // wechat, alipay, mock
    PayAmount     int64      `gorm:"not null"`
    Status        string     `gorm:"size:20;default:pending;index"` // pending, success, failed
    TransactionNo string     `gorm:"size:64;uniqueIndex"`           // 第三方支付流水号（幂等键）
    PaidAt        *time.Time
    NotifyRaw     string     `gorm:"type:text"`                     // 回调原始数据（对账用）
    NotifySignOk  bool       `gorm:"default:false"`                 // 签名验证结果
}

// OrderLog — 订单状态审计（不可变）
type OrderLog struct {
    BaseModel
    OrderID    uint   `gorm:"index:idx_order_logs;not null"`
    FromStatus string `gorm:"size:20"`                              // 空字符串表示初始创建
    ToStatus   string `gorm:"size:20;not null"`
    Operator   string `gorm:"size:32"`                              // "user:{id}", "system", "admin:{id}"
    Note       string `gorm:"size:256"`
}
```

### 订单状态机

```
                    ┌─────────────┐
                    │ pending_    │ ──(系统超时 30 分)──┐
                    │ payment     │                      │
                    └──────┬──────┘                      │
                      ┌────┴────┐                        │
                      │         │                        │
                      ▼         │                        ▼
              ┌───────────┐     │               ┌────────────┐
              │ pending_  │     │(用户取消)      │ cancelled  │
              │ delivery  │     │               │ (库存补偿)  │
              └─────┬─────┘     │               └────────────┘
                    │           │
                    ▼           │
              ┌───────────┐     │
              │ pending_  │     │
              │ receipt   │     │
              └─────┬─────┘     │
                    │           │
                    ▼           │
              ┌───────────┐     │
              │ completed │─────┘ (再次购买)
              └───────────┘
```

| 操作 | 原状态 | 目标状态 | 约束 | 补偿行为 |
|------|--------|---------|------|---------|
| 下单创建 | — | pending_payment | 乐观锁扣库存 | ROLLBACK 自动恢复 |
| 用户取消 | pending_payment | cancelled | 事务内补偿 | `stock += qty`, 写 OrderLog |
| 系统超时取消 | pending_payment | cancelled | 定时扫描 >30min | `stock += qty`, 写 OrderLog |
| 支付成功 | pending_payment | pending_delivery | 金额二次校验 | 幂等拦截 |
| 卖家发货 | pending_delivery | pending_receipt | 模拟物流单号 | — |
| 确认收货 | pending_receipt | completed | — | — |
| 再次购买 | completed | — | 将订单商品重新加入购物车 | — |

### 下单事务（核心流程）

```
┌──────────────────────────────────────────────────────────────────┐
│ TX BEGIN                                                         │
│  ① 幂等检查：查询 Idempotent-Key 是否已处理                      │
│     → 已处理：返回已有订单，不走后续逻辑                          │
│     → 未处理：记录 key → order_id 映射                           │
│  ② 生成 order_no：YYYYMMDDHHmmss + 4随机数字 + 2校验码 = 20位   │
│  ③ 创建 Order (status = pending_payment)                         │
│  ④ 遍历用户勾选的购物车商品：                                     │
│     a. 验证 product.status = 'on'，否则排除+收集提示              │
│     b. 验证 SKU 存在且激活                                        │
│     c. 从 DB 读取当前 price（绝不信任客户端）                     │
│     d. 乐观锁扣减：                                               │
│        UPDATE product_skus                                        │
│        SET stock = stock - ?, version = version + 1               │
│        WHERE id = ? AND stock >= ? AND version = ?                │
│     e. IF RowsAffected == 0 → ROLLBACK（库存不足或并发冲突）      │
│     f. 创建 OrderItem（price/name/image 快照）                    │
│     g. 累加 total_amount                                          │
│  ⑤ 服务端重算：pay_amount = total_amount + shipping_fee（微服务） │
│  ⑥ 更新 order.total_amount 和 order.pay_amount                   │
│  ⑦ 清空已购购物车商品（DELETE FROM cart WHERE user_id=? AND selected=1） │
│  ⑧ 写入 OrderLog（创建审计）                                      │
│ TX COMMIT                                                         │
│ 失败或 ROLLBACK → 返回 500/409，无副作用                           │
└──────────────────────────────────────────────────────────────────┘
```

### 第 7 周 — 购物车（7 个接口）

| 方法 | 路径 | 说明 | 安全要点 |
|------|------|------|---------|
| `GET` | `/api/v1/cart` | 购物车列表（含选中状态+汇总） | 仅返回当前用户 |
| `POST` | `/api/v1/cart` | 加入购物车 | qty ≤ stock（软提示），重复 SKU 合并数量 |
| `PUT` | `/api/v1/cart/:id` | 修改数量 | qty ≥ 1 |
| `DELETE` | `/api/v1/cart/:id` | 删除商品 | 验证归属 |
| `POST` | `/api/v1/cart/select` | 切换选中（单个或全选） | — |
| `POST` | `/api/v1/cart/merge` | 登录后合并本地购物车 | 本地→服务端合并；相同 SKU 取 max(qty) |
| `GET` | `/api/v1/cart/count` | 购物车数量（红点徽章） | 仅返回数量 |
| `GET` | `/api/v1/products/recommend` | "你可能还喜欢"推荐 | 当前分类随机取 |

**边界处理：**
- SKU 下架/删除 → 移出购物车+接口返回 `removed_items` 提示
- 价格变动 → 服务端当前价为准，前端不缓存价格
- 库存不足时加购 → qty 上限 = min(请求qty, stock)
- 购物车为空 → 返回空列表+推荐，不报错

**合并策略**（POST /cart/merge）：
```
请求体: [{ sku_id, quantity }]  // 本地购物车
处理逻辑:
  FOR each local_item:
    IF 服务端已存在相同 sku_id:
      new_qty = max(服务端.qty, 本地.qty)  // 取大
    ELSE:
      创建新 cart 记录
  返回合并后完整列表
```

### 第 8 周 — 订单（8 个接口）

| 方法 | 路径 | 说明 | 安全要点 |
|------|------|------|---------|
| `POST` | `/api/v1/orders` | **创建订单** | **Idempotent-Key 头**、事务、乐观锁扣库存、服务端算价 |
| `GET` | `/api/v1/orders` | 订单列表（按状态筛选+分页） | 仅当前用户、cursor 分页 |
| `GET` | `/api/v1/orders/:id` | 订单详情 | 验证 order.user_id == current_user |
| `POST` | `/api/v1/orders/:id/cancel` | **取消订单** | 仅 pending_payment、补偿库存、写审计 |
| `POST` | `/api/v1/orders/:id/confirm` | 确认收货 | 仅 pending_receipt |
| `POST` | `/api/v1/orders/:id/pay` | 模拟支付 | 金额二次校验、幂等 |
| `POST` | `/api/v1/orders/:id/logistics` | 查看物流（mock） | — |
| `POST` | `/api/v1/orders/:id/rebuy` | **再次购买** | 将订单商品加入购物车，跳过已下架 SKU |

**接口约束汇总：**

| 状态 | 可取消 | 可支付 | 可发货 | 可收货 | 可见 |
|------|--------|--------|--------|--------|------|
| pending_payment | ✅ | ✅ | ❌ | ❌ | 用户 |
| pending_delivery | ❌ | ❌ | ✅(管理员) | ❌ | 用户 |
| pending_receipt | ❌ | ❌ | ❌ | ✅ | 用户 |
| completed | ❌ | ❌ | ❌ | ❌ | 用户 |
| cancelled | ❌ | ❌ | ❌ | ❌ | 用户 |

**超时自动取消（系统补偿）：**
```
goroutine 定时器 (每分钟):
  SELECT * FROM orders
  WHERE status = 'pending_payment'
    AND created_at < NOW() - INTERVAL '30 minutes'
  FOR each order:
    TX BEGIN
      UPDATE orders SET status='cancelled', cancelled_at=NOW() WHERE id=? AND version=?
      FOR each order_item:
        UPDATE skus SET stock = stock + qty WHERE id = ?
      写入 OrderLog (operator = "system")
    TX COMMIT
  日志: "auto-cancelled N expired orders"
```

**再次购买流程（Rebuy）：**
```
FOR each order_item:
   检查 SKU 是否存在且 product.status == 'on'
   → 是：INSERT INTO cart (user_id, sku_id, quantity) ON CONFLICT (user_id, sku_id) DO UPDATE SET quantity = EXCLUDED.quantity ✅
   → 否：跳过，收集提示信息
 返回最终购物车列表
```

**测试清单（TDD 测试点）：**

下单接口：
- 正常下单 → 返回 order + 扣库存 + 清购物车 + 写审计
- 库存不足 → 409 错误，不回滚其他正常商品（整单失败）
- 同一 idempotent_key 重复请求 → 返回同一次订单
- 缺失 Idempotent-Key → 400
- SKU 不存在 → 400
- 购物车为空（所有商品不可用）→ 400
- 备注超长 → 400
- 金额从服务端读，客户端传入的价格被覆盖

状态流转：
- 下单 → pending_payment
- pending_payment → 取消 → cancelled + 库存恢复
- pending_payment → 支付 → pending_delivery
- pending_delivery → 发货 → pending_receipt
- pending_receipt → 确认 → completed
- 不可逆：已取消不能再支付，已完成不能再取消
- 不存在的订单 → 404
- 他人的订单 → 403

列表/详情：
- 按状态过滤分页
- 订单卡片包含商品快照（name/image/price/qty）
- 操作按钮根据状态返回

## 阶段五：支付与营销（TDD）

### 第 9 周 — 支付系统 ✅

> **架构原则**：先实现 Mock 支付，接口抽象为 `PaymentProvider`，真实微信/支付宝后期替换。
> Mock 支付模拟完整签名流程（HMAC-SHA256），使替换时仅新增 provider 实现，不改业务代码。

**支付抽象接口：**
```go
type Provider interface {
    CreatePayment(order *Order) (*Payment, error)
    ProcessCallback(notifyData []byte) (*CallbackResult, error)
}
```

**回调安全处理流程（核心）：**
```
收到支付回调 POST /api/v1/payment/callback/:method
  ① 读取原始 body（notify_raw）
  ② 验证签名（HMAC-SHA256）→ 失败则返回 200（不暴露信息），记录日志
  ③ 检查 transaction_no 是否已处理 → 已处理返回 200（幂等）
  ④ 查询 order，验证 order.pay_amount == callback.amount
     → 不匹配：记录安全告警日志，人工介入，返回 200
  ⑤ 验证 order.status == 'pending_payment'
     → 非待支付：记录异常日志，返回 200（幂等兜底）
  ⑥ 事务：
     a. 创建/更新 Payment：status=success, transaction_no, notify_sign_ok=true
     b. 更新 Order：status=pending_delivery, paid_at=NOW()
     c. 写入 OrderLog
  ⑦ 返回 200 给支付网关
```

| 方法 | 路径 | 说明 | 安全要点 |
|------|------|------|---------|
| `GET` | `/api/v1/payment/methods` | 可用支付方式列表 ✅ | 公开 |
| `POST` | `/api/v1/payment/pay` | 发起支付（选择方式+金额) ✅ | 金额二次确认 |
| `POST` | `/api/v1/payment/callback/:method` | **支付回调** ✅ | **签名验证、金额比对、幂等、IP白名单预留** |
| `GET` | `/api/v1/payment/status/:order_no` | 查询支付状态 ✅ | — |
| `POST` | `/api/v1/payment/refund` | **退款** ❌（待第 10 周） | 仅 super_admin/support 可操作 |

**Mock 支付细节：**
- ✅ 服务端生成 mock 交易号 `MOCK + YYYYMMDDHHmmss + 4随机`
- ✅ HMAC-SHA256 签名模拟
- ✅ 复用同一份 `processCallback` 逻辑
- ✅ 真实网关接入时仅新增 Provider 实现

### 第 10 周 — 营销系统 ✅

**优惠券模型：**
```go
type Coupon struct {
    BaseModel
    Name         string    `gorm:"size:64;not null"`       // 优惠券名称
    Type         string    `gorm:"size:20;not null"`       // full_reduce(满减), discount(折扣)
    Condition    int64     `gorm:"not null"`               // 满减门槛（分）
    Discount     int64     `gorm:"not null"`               // 减免金额（分）或折扣比例
    TotalCount   int       `gorm:"default:0"`              // 发行总量（0=不限）
    RemainCount  int       `gorm:"default:0"`              // 剩余数量
    PerLimit     int       `gorm:"default:1"`              // 每人限领
    StartAt      time.Time `gorm:"not null"`               // 有效期开始
    EndAt        time.Time `gorm:"not null"`               // 有效期结束
    Status       string    `gorm:"size:20;default:active"` // active, expired, disabled
}

type UserCoupon struct {
    BaseModel
    UserID     uint       `gorm:"index;not null"`
    CouponID   uint       `gorm:"index;not null"`
    UsedAt     *time.Time
    OrderID    *uint      `gorm:"index"`                   // 使用时关联订单
    Status     string     `gorm:"size:20;default:unused"`  // unused, used, expired
}
```

| 方法 | 路径 | 说明 | 安全要点 |
|------|------|------|---------|
| `POST` | `/api/v1/admin/coupons` | 创建优惠券 ✅ | super_admin/operator |
| `POST` | `/api/v1/coupons/:id/receive` | 领取优惠券 ✅ | 限领检查、剩余量检查 |
| `GET` | `/api/v1/coupons/available` | 结算页可用优惠券 ✅ | 按订单金额筛选，返回可用的 |
| `POST` | `/api/v1/coupons/calculate` | 计算优惠后金额 ✅ | 服务端计算，不信任客户端 |

**限时秒杀：**
- ✅ FlashSale 模型：`product_id, sku_id, flash_price, flash_stock, start_at, end_at, status`
- ✅ 库存独立于普通库存，扣减使用乐观锁
- ✅ 秒杀接口返回倒计时（秒）
- GET `/api/v1/flash-sales` 公开接口

**积分系统：**
- ✅ `PointLog`流水表（`point_logs`），记录积分变动
- ✅ 支付成功自动赠送积分（1 积分/100 分消费金额）
- ✅ `GET /api/v1/points` 查询余额
- ✅ `GET /api/v1/points/logs` 积分流水
- ❌ 积分抵扣（待订单结算页集成）

### 前期 Bug 修复（第 8 周后补充）

- [x] **AutoMigrate 缺失** — main.go 启动后未自动建表
- [x] **Viper 环境变量不生效** — 添加 `SetEnvKeyReplacer` 使 `GOSH_*` 环境变量可覆盖配置
- [x] **Zap Logger 配置未使用** — 从配置文件读取日志级别和文件路径
- [x] **无优雅关闭** — 添加信号处理 + Scheduler 停止通道
- [x] **Rebuy 数量策略错误** — 从累加改为替换（匹配 roadmap 规范）

## 阶段六：部署与优化（第 11-12 周）

- [x] Docker 容器化（多阶段构建）
- [x] docker-compose（后端 + PG + Redis）
- [x] 集成测试覆盖（Python 冒烟测试 + k6 负载测试，24/24 通过）
- [x] Redis 缓存（热点数据、分布式锁、限流中间件、购物车）
  - [x] 分类树、首页轮播、品牌故事、热搜排行缓存
  - [x] 商品详情、秒杀活动列表缓存
  - [x] 购物车列表缓存（每用户 5min TTL，写操作失效）
  - [x] 分布式锁（pkg/cache/lock.go）
  - [x] 滑动窗口限流中间件（基于用户/IP + 路径，按角色分级）
  - [x] Cache-Aside 模式，Redis 不可用时自动降级
- [ ] CI/CD 流水线（GitHub Actions）
  - [ ] PR 自动运行 `go test ./...`
  - [ ] 合并到 main 自动部署
- [x] 消息队列（RabbitMQ，异步积分赠送 + 支付回调处理）
  - [x] pkg/mq — 连接管理、Exchange/Queue 声明、发布者、消费者框架
  - [x] internal/worker — PointWorker（订单支付后异步送积分）、PaymentWorker（支付回调异步处理）
  - [x] 支付成功有 MQ 时走 MQ，无 MQ 时降级为同步送积分
  - [x] docker-compose 集成 RabbitMQ 4
- [x] 性能优化（索引、N+1 查询修复、搜索优化）
  - [x] 添加缺失索引：Order.created_at、OrderItem.SKUID、FlashSale(status,start_at,end_at)、Coupon(status,start_at,end_at)、Address(user_id,is_default)、User.status、Product.{is_new,is_hot,is_featured}
  - [x] 修复购物车 N+1（2N+1 → 3 次查询：批量 IN 加载 SKU + Product）
  - [x] 修复 Rebuy N+1（3N+1 → 批量加载 SKU + Product）
  - [x] isProductOnline 轻量化（仅查 status 列，不再读取全行）
  - [x] 搜索优化：PostgreSQL 使用 to_tsvector FTS，SQLite 保持 LIKE 兼容
  - [x] User List 添加 ORDER BY id asc
  - [x] Coupon List 添加 LIMIT 100 防止全表扫描
- [ ] API 文档自动生成（Swagger）
- [ ] 线上部署

## 阶段七：服务端细节打磨（第 13 周）

> 阶段六的基于 Docker 的综合测试发现了若干设计瑕疵和遗漏的防护，
> 在启动多端前端之前集中修补，避免多端联调时重复返工。

### 7.1 DTO 设计修正 ✅

| 问题 | 现状 | 目标 |
|------|------|------|
| `Images` 字段类型 | `string`（存 JSON 序列化字符串） | `[]string` + `gorm:"type:text;serializer:json"` ✅ |
| 订单 `CreateOrderRequest` 无 items | 只能从购物车下单 | 支持 `items: []CreateOrderItem`，无 items 时回退购物车 ✅ |
| 购物车 `select` 字段名与 JS 关键字冲突 | `json:"select"` | 改为 `json:"selected"` ✅ |
| `CalculateCouponRequest` 需预知 coupon_id | 无法做"自动最优" | `coupon_id` 改为 `*uint`，不传时自动选最优 ✅ |

### 7.2 安全与输入防护增强 ✅

- [x] **购物车数量上限**：`AddCartRequest` / `UpdateCartRequest` 补充 `max=99` binding 校验
- [x] **上传文件大小**：Gin 全局 `MaxMultipartMemory = 10MB`，Upload handler 内额外校验
- [x] **XSS 防护**：`middleware.SanitizeInput()` 全局中间件，拦截 POST/PUT/PATCH 对 JSON body 内所有字符串做 `html.EscapeString`
> 分页默认值已存在于所有 endpoint（page=1, size=10 默认值 + 负值/零值 400 校验），此项无需额外工作。

### 7.3 订单与库存策略优化 ✅

- [x] **库存扣减时机评估**：保持当前单阶段模型（下单即扣库存，乐观锁，超时释放）。`reserved_stock` 分离待后续秒杀高并发场景再引入。
- [x] **订单超时取消**：改为配置化，`config.yaml` → `order.timeout_minutes`，Scheduler 启动时读取。
- [x] **Rebuy 补充**：返回结构化 `RebuyResponse{Cart, SkippedItems}`，已下架商品明确提示原因而非 500。

### 7.4 API 契约一致性 ✅

- [x] **统一错误码**：`pkg/errcode` 模块化错误码（1xxx-6xxx），全 handler 使用 `*WithCode` 变体。
- [x] **幂等性响应**：重复请求返回完整的 `OrderResponse`（存于 `IdempotencyRecord.Response`）。
- [x] **购物车库存实时性**：列表返回 `MaxBuyable` 字段（`min(stock, 99)`）。

### 7.5 可观测性 ✅

- [x] **请求追踪**：`middleware.RequestID()` 为每个请求注入唯一 `request_id`，贯穿日志和 `X-Request-ID` 响应头。
- [x] **结构化错误日志**：Logger 中间件按状态分级：`2xx`→info，`4xx`→warn，`5xx`→error，日志包含 `request_id`、`user_id`、`path`、`method`、`cost`。
- [x] **慢查询日志**：GORM 自定义 logger，`SlowThreshold: 200ms`，生产环境默认 `Warn` 级别。

### 7.6 已修复

- [x] **搜索 SQL 注入**：`to_tsquery` → `plainto_tsquery` + `sanitizeFTS` 清洗输入，单引号等特殊字符不再触发 500（`server/internal/repository/product/product.go:147`）

## 阶段八：多客户端前端（待定）

> **⚠️ 前端开发暂缓。** 后端 API 已全部就绪（89 endpoint，Swagger 文档完整），待前端资源到位后启动。
> 以下四端可并行开发，共享同一套 API 契约。
> 设计稿来源：阶段一的 uniapp 原型 + 各端补充设计。

### 8.1 C端 — 现有 UniApp 重构（普通用户）

| 子任务 | 说明 |
|--------|------|
| 原型标注与页面联调 | 将原型 HTML 落地为 Vue 页面，对接真实 API |
| 商品详情页 | product-detail.html → /pages/product/detail |
| 搜索结果页 | search.html → /pages/search/search |
| 结算页 | checkout.html → /pages/order/checkout |
| 订单列表/详情 | orders.html → /pages/order/list + /pages/order/detail |
| 现有页面 API 对接 | 首页、分类、购物车、个人中心替换 Mock 数据 |

### 8.2 B端 — Admin Web（超级管理员）

| 子任务 | 说明 |
|--------|------|
| 技术选型 | 推荐 Vue 3 + Element Plus / Ant Design |
| 仪表盘 | 销售概览、用户统计、商品统计 |
| 用户管理 | 用户列表、角色分配、商户审核 |
| 商品管理 | 商品 CRUD、分类管理、轮播/品牌管理 |
| 订单管理 | 订单列表/详情、发货操作、退款处理 |
| 营销管理 | 优惠券创建、秒杀设置 |

### 8.3 运营端 — UniApp（运营 + 客服）

| 子任务 | 说明 |
|--------|------|
| 独立 UniApp 项目 | 复用 C 端组件库，独立 TabBar 和路由 |
| 商品审核 | 商品上架/下架、编辑 |
| 售后处理 | 退款/退货审核 |
| 用户咨询 | 客服对话功能 |
| 数据看板 | 运营数据概览 |

### 8.4 商家端 — UniApp（商户）

| 子任务 | 说明 |
|--------|------|
| 独立 UniApp 项目 | 面向入驻商户，仅可见自身店铺数据 |
| 店铺管理 | 店铺信息编辑 |
| 商品管理 | 自有商品 CRUD（受 `tenant_id` 隔离） |
| 订单管理 | 自有订单处理、发货 |
| 数据统计 | 店铺销售数据 |
