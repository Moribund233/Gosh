# Gosh Mall 外部测试报告

> 日期: 2026-06-15
> 测试环境: Linux, SQLite (dev), Go 1.x, Gin

---

## 1. Python 冒烟测试 (api_test.py)

**结果: 24/24 ✅ 全部通过**

| # | 测试项 | 状态 | 说明 |
|---|--------|------|------|
| 1 | Health check | ✅ | `GET /health` → 200 |
| 2 | User register | ✅ | `POST /user/register` → 201, 返回 token |
| 3 | Token returned | ✅ | data.token 非空 |
| 4 | Get profile | ✅ | `GET /user/profile` → 200 |
| 5 | Profile nickname | ✅ | nickname == "测试用户" |
| 6 | Create address | ✅ | `POST /addresses` → 201 |
| 7 | Get categories | ✅ | `GET /categories` → 200 |
| 8 | Get products (empty) | ✅ | `GET /products` → 200, data.list 空 |
| 9 | Empty list | ✅ | len(list) == 0 |
| 10 | Search products | ✅ | `GET /products/search?keyword=test` → 200 |
| 11 | Get banners | ✅ | `GET /banners` → 200 |
| 12 | Payment methods | ✅ | `GET /payment/methods` → 200 |
| 13 | 3 methods | ✅ | mock/wechat/alipay |
| 14 | Flash sales | ✅ | `GET /flash-sales` → 200 |
| 15 | Query points | ✅ | `GET /points` → 200 |
| 16 | Initial points = 0 | ✅ | data.points == 0 |
| 17 | Point logs | ✅ | `GET /points/logs` → 200 |
| 18 | Favorites | ✅ | `GET /favorites` → 200 |
| 19 | Browse history | ✅ | `GET /browse-history` → 200 |
| 20 | Cart | ✅ | `GET /cart` → 200 |
| 21 | Cart count | ✅ | `GET /cart/count` → 200 |
| 22 | Orders | ✅ | `GET /orders` → 200 |
| 23 | Available coupons | ✅ | `GET /coupons/available` → 200 |
| 24 | Unauthorized rejected | ✅ | 无 token 访问认证接口 → 401 |

---

## 2. 并发压力测试 (Python Threading)

**场景:** 20 并发用户同时注册 + 获取 Profile

| 指标 | 值 |
|------|----|
| 总请求数 | 40 (20 register + 20 profile) |
| 通过 | 40 ✅ |
| 失败 | 0 |
| 平均响应时间 | 0.316s |
| 最小响应时间 | 0.187s |
| 最大响应时间 | 0.442s |
| p95 | 0.442s |

> 注: SQLite 单写者模式下达到此性能，生产切换 PostgreSQL 后并发能力会大幅提升。

---

## 3. k6 压测 (Docker + PostgreSQL)

**场景:** 5 并发 VU，50s 阶梯增压，模拟完整用户流（注册 → Profile → 各类查询 → 创建地址）

| 指标 | 值 |
|------|----|
| 总检查数 | 2280 |
| 通过率 | 100% ✅ |
| 失败 | 0 |
| 总请求数 | 2280 |
| 平均响应时间 | 7.26ms |
| p95 | 71.46ms |
| 错误率 | 0.00% |
| 迭代数 | 190 (3.79/s) |

**所有 12 项检查通过:** health, banners, categories, products, payment methods, register, profile, points, cart, orders, favorites, address

---

## 4. 容器化配置

```
# 构建 & 运行（独立容器，连接已有 PostgreSQL）
docker build -t gosh-server:latest .
docker run -d --name gosh-server --network 1panel-network \
  -p 9292:8080 \
  -e GOSH_DATABASE_DRIVER=postgres \
  -e GOSH_DATABASE_HOST=postgresql \
  -e GOSH_DATABASE_PORT=5432 \
  -e GOSH_DATABASE_USER=alucard \
  -e GOSH_DATABASE_PASSWORD=admin123456 \
  -e GOSH_DATABASE_DBNAME=gosh \
  -e GOSH_DATABASE_SSLMODE=disable \
  -e GIN_MODE=release \
  gosh-server:latest

# 或使用 docker-compose（自带 PostgreSQL）
docker compose up -d
```

**多阶段构建:** builder 阶段使用 `golang:1.26-alpine` + `GOPROXY=https://goproxy.cn,direct`；runtime 阶段使用 `alpine:latest`，仅含 ca-certificates + tzdata。

---

## 5. 路由覆盖率验证

启动时 Gin 注册了 **全部 63 条路由**，确认新功能路由已生效:

| 新增路由 | 方法 | 说明 |
|----------|------|------|
| `POST /api/v1/orders/:id/apply-points` | ✅ | 积分抵扣 |
| `POST /api/v1/admin/payment/refund` | ✅ | 退款处理 |

---

## 6. 已知问题（已解决）

| 问题 | 解决方案 |
|------|----------|
| k6 硬编码 port 8080 | 改为 `__ENV.BASE_URL` 环境变量，默认 9292 |
| k6 手机号生成位数错误 | `Date.now().slice(-10)` → `slice(-8)`，保证 11 位 |
| Viper 环境变量不生效 | `config.yaml` 缺少 PG 字段（`host`, `port`, `user` 等），`AllSettings()` 不会枚举未知 key，导致 `Unmarshal` 忽略。修复：在 yaml 中声明所有 PG key |

---

## 7. 后续优化方向

- **静态文件服务** — `/uploads/` 目录需挂载 volume
- **Swagger 文档** — 接入 swaggo 自动生成 API 文档
- **Redis 缓存** — 热点数据（banner, category, product）加缓存
- **消息队列** — 订单超时、异步通知场景
- **CI/CD** — GitHub Actions 自动构建 + 部署

---

## 8. 结论

**全部 24 项 Python 冒烟测试 + k6 压测（2280 checks, 100%）均通过。Docker 多阶段构建 + PostgreSQL 环境就绪。**
