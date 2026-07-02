# 谷穗 C 端（uni-app）开发规划

> **品牌**: 谷穗 | **框架**: uni-app 3 + Vue 3 + TypeScript + Vite 5
> **组件库**: ThorUI | **状态管理**: Pinia | **后端**: 57 个 C 端接口已就绪

---

## 一、技术选型

| 类别 | 选型 | 说明 |
|------|------|------|
| 组件库 | **ThorUI** | `npm install thorui-uni`，easycom 自动导入，前缀 `tui-` |
| 图标 | ThorUI 内置 `tui-icon` + PNG 图片 | 禁止硬编码 SVG 和 emoji 图标 |
| HTTP 请求 | 自定义 `request.ts`（已就绪）+ ThorUI `tui.request` 备用 | 统一鉴权、错误处理 |
| 状态管理 | Pinia（已安装，待创建 Store） | auth / cart / user / product |
| 路由 | uni-app 原生 `pages.json`（已定义 13 条路由） | |

---

## 二、后端 API 映射

57 个 C 端接口，按优先级分组：

| 域 | 接口数 | 依赖 Phase |
|----|--------|------------|
| 认证（注册/登录） | 2 | Phase 1 |
| 用户信息 | 2 | Phase 1 |
| Banner / 品牌故事 | 2 | Phase 2 |
| 分类树 | 2 | Phase 2 |
| 商品列表/搜索/详情 | 6 | Phase 2-3 |
| 评价 | 2 | Phase 3 |
| 秒杀 | 1 | Phase 2 |
| 购物车（7个） | 7 | Phase 4 |
| 订单（8个） | 8 | Phase 4 |
| 支付 | 2 | Phase 4 |
| 地址 CRUD | 4 | Phase 4 |
| 收藏 / 浏览记录 | 5 | Phase 5 |
| 优惠券 / 积分 | 5 | Phase 5 |

---

## 三、阶段规划与进度追踪

### Phase 1 — 基础设施

| # | 任务 | 文件 | 状态 |
|---|------|------|------|
| 1.1 | 安装 ThorUI（npm + easycom 配置） | `pages.json`, `package.json` | ✅ |
| 1.2 | 配置环境变量（dev/prod BASE_URL） | `vite.config.ts`, `.env` | ✅ |
| 1.3 | 创建 Auth Store（token 管理、登录态） | `src/stores/auth.ts` | ✅ |
| 1.4 | 创建 User Store（用户信息） | `src/stores/user.ts` | ✅ |
| 1.5 | 创建 Cart Store（购物车状态） | `src/stores/cart.ts` | ✅ |
| 1.6 | 登录页 UI + API 对接 | `pages/login/login.vue` | ✅ |
| 1.7 | 注册页 UI + API 对接 | `pages/login/register.vue` | ✅ |
| 1.8 | 全局导航守卫（未登录拦截 + Token 过期处理） | `App.vue` | ✅ |
| 1.9 | 安装 ThorUI `tui-icon`，替换 tabbar PNG 图标 | `pages.json` | ✅ |

### Phase 2 — 首页 & 发现

| # | 任务 | 文件 | 状态 |
|---|------|------|------|
| 2.1 | 首页 — 轮播 Banner | `pages/index/index.vue` | ✅ |
| 2.2 | 首页 — 品牌故事卡片 | `pages/index/index.vue` | ✅ |
| 2.3 | 首页 — 快捷菜单（4项入口） | `pages/index/index.vue` | ✅ |
| 2.4 | 首页 — 秒杀倒计时 + 横向滚动 | `pages/index/index.vue` | ✅ |
| 2.5 | 首页 — 商品瀑布流（2列网格） | `pages/index/index.vue` | ✅ |
| 2.6 | 分类页 — 左侧分类树 | `pages/category/category.vue` | ✅ |
| 2.7 | 分类页 — 右侧子分类网格 + Banner | `pages/category/category.vue` | ✅ |
| 2.8 | 搜索页 — 搜索栏 | `pages/search/search.vue` | ✅ |
| 2.9 | 搜索页 — 历史搜索 + 热门搜索 | `pages/search/search.vue` | ✅ |
| 2.10 | 搜索页 — 搜索结果列表 | `pages/search/search.vue` | ✅ |

### Phase 3 — 商品详情

| # | 任务 | 文件 | 状态 |
|---|------|------|------|
| 3.1 | 商品详情 — 图片画廊 | `pages/sub-package-product/product/product.vue` | ✅ |
| 3.2 | 商品详情 — 价格/标题/标签 | `pages/sub-package-product/product/product.vue` | ✅ |
| 3.3 | 商品详情 — SKU 规格选择 | `pages/sub-package-product/product/product.vue` | ✅ |
| 3.4 | 商品详情 — 促销信息栏 | `pages/sub-package-product/product/product.vue` | ✅ |
| 3.5 | 商品详情 — 评价摘要 | `pages/sub-package-product/product/product.vue` | ✅ |
| 3.6 | 商品详情 — 底部操作栏（收藏/加购/立即购买） | `pages/sub-package-product/product/product.vue` | ✅ |

### Phase 4 — 购物车 & 结算

| # | 任务 | 文件 | 状态 |
|---|------|------|------|
| 4.1 | 购物车 — 列表（勾选/数量/价格） | `pages/sub-package-order/cart/cart.vue` | ✅ |
| 4.2 | 购物车 — 推荐商品区 | `pages/sub-package-order/cart/cart.vue` | ✅ |
| 4.3 | 购物车 — 底部结算栏 | `pages/sub-package-order/cart/cart.vue` | ✅ |
| 4.4 | 地址管理 — 列表 | `pages/address/address.vue` | ✅ |
| 4.5 | 地址管理 — 新增/编辑 | `pages/address-edit/address-edit.vue` | ✅ |
| 4.6 | 结算页 — 地址选择 | `pages/sub-package-order/checkout/checkout.vue` | ✅ |
| 4.7 | 结算页 — 商品清单 | `pages/sub-package-order/checkout/checkout.vue` | ✅ |
| 4.8 | 结算页 — 优惠券选择 | `pages/sub-package-order/checkout/checkout.vue` | ✅ |
| 4.9 | 结算页 — 支付方式选择 | `pages/sub-package-order/checkout/checkout.vue` | ✅ |
| 4.10 | 结算页 — 订单备注 + 提交 | `pages/sub-package-order/checkout/checkout.vue` | ✅ |
| 4.11 | 订单创建（Idempotent-Key） | `services/order.ts` | ✅ |

### Phase 5 — 用户中心

| # | 任务 | 文件 | 状态 |
|---|------|------|------|
| 5.1 | 个人中心 — 用户信息头 | `pages/profile/profile.vue` | ✅ |
| 5.2 | 个人中心 — 订单统计（待付款/待发货/待收货/已完成/售后） | `pages/profile/profile.vue` | ✅ |
| 5.3 | 个人中心 — 工具栏（地址/优惠券/收藏/浏览记录） | `pages/profile/profile.vue` | ✅ |
| 5.4 | 个人中心 — 菜单组（会员/设置/关于） | `pages/profile/profile.vue` | ✅ |
| 5.5 | 订单列表 — 状态 Tab 切换 | `pages/sub-package-order/orders/orders.vue` | ✅ |
| 5.6 | 订单列表 — 订单卡片（状态/商品/操作按钮） | `pages/sub-package-order/orders/orders.vue` | ✅ |
| 5.7 | 订单操作（取消/付款/确认收货/再次购买） | `pages/sub-package-order/orders/orders.vue` | ✅ |
| 5.8 | 优惠券页面 | `pages/coupon/coupon.vue` | ✅ |
| 5.9 | 积分页面 | `pages/points/points.vue` | ✅ |
| 5.10 | 收藏/浏览记录（占位） | `pages/profile/profile.vue` | ⬜ |

---

## 四、开发规范

### 4.1 组件与图标

- **禁止** 在模板中硬编码 `<svg>` / `viewBox` / `<path d="...">`
- **禁止** 在 CSS 或 JS 中使用 `data:image/svg+xml`
- **禁止** 在代码中使用 emoji 字符
- 图标统一使用 ThorUI `tui-icon` 组件或 PNG 图片资源

### 4.2 目录结构

```
src/
├── components/          # 公共组件（按需创建）
│   ├── product-card/    # 商品卡片
│   ├── order-card/      # 订单卡片
│   └── ...
├── stores/              # Pinia 状态
│   ├── auth.ts
│   ├── cart.ts
│   └── user.ts
├── services/            # API 服务（已就绪）
├── pages/               # 页面（已定义路由）
├── utils/               # 工具函数（已就绪）
└── static/              # 静态资源
```

### 4.3 API 对接约定

- 所有请求走 `services/request.ts`（已封装 token 注入 + 401 自动跳转登录）
- 新建 service 模块时按域名拆分（如 `services/order.ts`）
- 优先使用 TypeScript 接口类型（已有 Product/Category/Banner 等类型定义）

---

## 五、原型参考

- 谷穗原型位于 `prototype/gusui/`
- 覆盖页面：首页、分类、商品详情、购物车、结算、订单、个人中心、搜索

---

## 六、状态图例

| 符号 | 含义 |
|------|------|
| ⬜ | 未开始 |
| 🔄 | 进行中 |
| ✅ | 已完成 |
| ❌ | 已取消/阻塞 |
