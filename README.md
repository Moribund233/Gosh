# Gosh Mall

基于 Go (Gin) + Vue (UniApp) 的全栈电商商城系统，涵盖用户认证、商品管理、购物车、订单、支付集成、优惠券、秒杀、积分等完整电商功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| **后端** | Go 1.26, Gin, GORM, Viper, JWT (golang-jwt) |
| **数据库** | PostgreSQL / SQLite（GORM 自动迁移） |
| **缓存** | Redis（旁路缓存、分布式锁、限流） |
| **消息队列** | RabbitMQ（异步积分赠送、支付处理） |
| **前端** | UniApp (Vue 3 + TypeScript + Pinia + uView Plus) |
| **认证** | JWT（基于角色：super_admin, operator, merchant, support, user） |
| **基础设施** | Docker, docker-compose, k6 |
| **测试** | testify, go-sqlmock, Gin 测试工具 |
| **文档** | Swagger/OpenAPI |

## 架构

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   UniApp     │────▶│   Gin HTTP   │────▶│   Service    │
│  (Vue 3)     │     │   Handler    │     │   (业务逻辑)  │
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                 │
                                          ┌──────▼───────┐
                                          │  Repository  │
                                          │  (数据访问)   │
                                          └──────┬───────┘
                                                 │
                          ┌──────────────────────┼──────────────────────┐
                          │                      │                      │
                    ┌─────▼─────┐          ┌─────▼─────┐          ┌────▼────┐
                    │ PostgreSQL │          │   Redis   │          │ RabbitMQ│
                    │  / SQLite  │          │  (缓存)   │          │  (MQ)   │
                    └───────────┘          └───────────┘          └─────────┘
```

**分层设计，接口隔离，便于测试：**

- **Handler** — HTTP 层：请求校验、响应格式化，不含业务逻辑
- **Service** — 业务逻辑层：纯 Go 接口，测试时轻松 Mock
- **Repository** — 数据访问层：GORM 查询封装在接口之后

## 功能特性

### 用户系统
- 注册/登录，bcrypt 密码加密
- JWT 认证 + 基于角色的访问控制（RBAC）
- 个人信息管理（昵称、头像）
- 收货地址管理（CRUD、默认地址）
- 浏览记录 & 商品收藏
- 商户入驻申请与审核流程

### 商品系统
- 多级分类树 + 分类横幅
- SPU/SKU 模型，支持多规格
- 商品搜索（PostgreSQL 全文搜索 / SQLite LIKE）
- 商品评价与评分
- 热搜排行 & 用户搜索历史
- 图片上传（multipart + base64）

### 购物车与订单
- 购物车 CRUD，支持选中状态
- 登录后购物车合并（本地 → 服务端）
- **幂等订单创建**（Idempotent-Key 请求头）
- **乐观锁扣库存**（防超卖）
- **数据库事务**保证订单生命周期
- 订单状态机：pending_payment → pending_delivery → pending_receipt → completed
- 超时自动取消（30 分钟，补偿恢复库存）
- 再次购买：一键复购，校验库存，跳过下架商品
- 完整订单审计日志（OrderLog）

### 支付系统
- 可插拔 `PaymentProvider` 接口（Mock / 微信 / 支付宝）
- Mock 支付模拟 HMAC-SHA256 签名流程
- 安全回调处理：签名验证、金额二次校验、幂等
- 支付状态查询

### 营销
- 优惠券：满减 / 折扣类型，定时活动，每人限领
- 限时秒杀：独立秒杀库存，倒计时，乐观锁
- 积分系统：支付送积分（每 100 分消费送 1 积分），积分流水日志

### 管理后台
- 角色管理（super_admin, operator, 等）
- 用户管理（列表、角色分配、商户审核）
- 商品管理（CRUD、分类、轮播、品牌故事）
- 订单管理（列表、发货、退款）
- 优惠券创建、秒杀设置

### 基础设施
- Redis 缓存：分类树、轮播图、商品详情、秒杀、热搜（Cache-Aside，宕机自动降级）
- 分布式锁（pkg/cache/lock.go）
- 滑动窗口限流中间件（按角色区分：用户/IP + 路径）
- RabbitMQ 消息队列：异步积分赠送 + 支付处理
- 优雅关闭（信号处理）
- 请求 ID 追踪（贯穿日志）
- 结构化日志（zap），按级别分级
- 慢查询检测（GORM，200ms 阈值）

### 测试
- **TDD 驱动开发**：先写测试，再写实现
- 单元测试：go-sqlmock + testify
- HTTP Mock 测试工具
- Python 端到端冒烟测试（24/24 通过）
- k6 负载测试脚本

## 快速开始

### 开发环境（SQLite）

```bash
cd server
go run cmd/server/main.go
```

服务启动于 `http://localhost:8080`，Swagger 文档访问 `http://localhost:8080/swagger/index.html`。

### 生产环境（PostgreSQL + Docker）

```bash
cd server
docker compose up -d
```

服务启动于 `http://localhost:9300`。

### 配置

配置文件：`server/config/config.yaml`，所有配置项均可通过 `GOSH_*` 环境变量覆盖：

```bash
GOSH_DATABASE_DRIVER=postgres \
GOSH_DATABASE_HOST=127.0.0.1 \
GOSH_DATABASE_PORT=5432 \
GOSH_DATABASE_USER=gosh \
GOSH_DATABASE_PASSWORD=gosh \
GOSH_DATABASE_DBNAME=gosh \
GOSH_DATABASE_SSLMODE=disable \
go run cmd/server/main.go
```

完整配置参考：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| server.port | 8080 | HTTP 监听端口 |
| server.mode | debug | Gin 模式（debug/release） |
| database.driver | sqlite | 数据库驱动（postgres / sqlite） |
| database.path | gosh.db | SQLite 文件路径 |
| database.host | - | PostgreSQL 地址 |
| database.port | - | PostgreSQL 端口 |
| database.user | - | PostgreSQL 用户名 |
| database.password | - | PostgreSQL 密码 |
| database.dbname | - | PostgreSQL 数据库名 |
| redis.host | 127.0.0.1 | Redis 地址 |
| redis.port | 6379 | Redis 端口 |
| rabbitmq.host | 127.0.0.1 | RabbitMQ 地址 |
| rabbitmq.port | 5672 | RabbitMQ 端口 |
| jwt.secret | gosh-dev-secret | JWT 签名密钥 |
| jwt.expire_hour | 72 | JWT Token 过期时间（小时） |
| upload.dir | storage/upload | 文件上传目录 |
| upload.max_size | 10 | 上传大小限制（MB） |
| order.timeout_minutes | 30 | 订单超时自动取消时间（分钟） |

## API 概览

总计 89 个端点，按业务域划分。运行时可访问 `/swagger/index.html` 查看 Swagger 文档。

| 域 | 端点 | 认证 |
|-----|--------|------|
| 用户 | 注册、登录、个人信息、更新 | 公开 / JWT |
| 地址 | CRUD | JWT |
| 收藏 | 添加、列表、删除 | JWT |
| 浏览记录 | 列表、清除 | JWT |
| 商户 | 申请、审核、列表 | JWT / 管理员 |
| 分类 | 树形结构、列表 | 公开 |
| 商品 | 列表、详情、搜索、推荐、管理 | 公开 / 管理员 |
| 评价 | 创建、列表 | JWT |
| 轮播 | 列表、管理 | 公开 / 管理员 |
| 上传 | 图片上传 | JWT |
| 购物车 | CRUD、选中、合并、数量 | JWT |
| 订单 | 创建、列表、详情、取消、确认、支付、物流、再次购买 | JWT |
| 支付 | 方式列表、发起、回调、状态、退款 | JWT / 公开回调 |
| 优惠券 | 创建、领取、可用、计算优惠 | 管理员 / JWT |
| 秒杀 | 活动列表 | 公开 |
| 积分 | 余额、流水 | JWT |

## 项目结构

```
├── LICENSE
├── README.md
├── docs/
│   ├── roadmap.md           # 开发路线图
│   └── external-test-report.md
├── server/                  # Go 后端
│   ├── cmd/server/main.go   # 入口文件
│   ├── config/              # YAML 配置文件
│   ├── internal/
│   │   ├── config/          # 配置加载（Viper）
│   │   ├── database/        # 数据库初始化
│   │   ├── handler/         # HTTP 处理器
│   │   ├── middleware/      # 鉴权、CORS、日志、恢复、限流、XSS 过滤
│   │   ├── model/           # GORM 数据模型
│   │   ├── dto/request/     # 请求 DTO
│   │   ├── dto/response/    # 响应 DTO
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑层
│   │   ├── router/          # 路由注册
│   │   ├── scheduler/       # 后台任务（订单超时取消）
│   │   └── worker/          # MQ 消费者（积分、支付）
│   ├── pkg/
│   │   ├── auth/            # JWT 工具
│   │   ├── cache/           # Redis 缓存 + 分布式锁
│   │   ├── errcode/         # 错误码（1xxx-6xxx）
│   │   ├── mq/              # RabbitMQ 发布/订阅
│   │   └── response/        # 统一响应格式
│   ├── scripts/             # 测试脚本（Python + k6）
│   ├── Dockerfile           # 多阶段构建
│   └── docker-compose.yml   # 服务编排（Server + PostgreSQL + Redis + RabbitMQ）
└── uniapp/                  # UniApp 前端（Vue 3 + TypeScript）
    ├── src/
    │   ├── App.vue
    │   ├── main.js
    │   ├── pages/           # 购物车、分类、首页、个人中心
    │   ├── stores/          # Pinia 状态管理
    │   └── utils/           # HTTP 请求封装
    └── vite.config.js
```

## 测试

```bash
# 单元测试（TDD 全覆盖）
cd server && go test ./...

# 端到端冒烟测试（需先启动服务）
pip install requests
python scripts/api_test.py

# 压力测试（需先启动服务）
k6 run -e BASE_URL=http://localhost:9300/api/v1 scripts/k6_test.js
```

## 开发路线图

详见 [docs/roadmap.md](docs/roadmap.md)。

## 开源许可

MIT
