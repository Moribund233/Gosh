# Gosh Mall

A full-stack e-commerce mall built with Go (Gin) + Vue (UniApp).

## Tech Stack

**Backend:** Go 1.26, Gin, GORM, Viper, JWT, PostgreSQL/SQLite

**Frontend:** UniApp (Vue 3 + TypeScript)

**Infra:** Docker, docker-compose, k6

## Architecture

```
Handler → Service(interface) → Repository(interface)
                                  │
                            GORM + PostgreSQL/SQLite
```

- **Handler** — HTTP layer, request validation, response formatting
- **Service** — Business logic, interfaces for testability
- **Repository** — Data access, interfaces for testable queries

## Quick Start

```bash
# Development (SQLite)
cd server
go run cmd/server/main.go

# Production (PostgreSQL + Docker)
cd server
docker compose up -d
```

Server starts at `http://localhost:8080` (dev) or `http://localhost:9292` (Docker).

### Configuration

Config is in `server/config/config.yaml` with overrides via `GOSH_*` env vars:

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

### Testing

```bash
# Unit tests
cd server && go test ./...

# Smoke tests (server must be running)
pip install requests
python scripts/api_test.py

# Load tests (server must be running)
k6 run -e BASE_URL=http://localhost:9292/api/v1 scripts/k6_test.js
```

## Project Structure

```
docs/                        # Docs & plans
server/
  cmd/server/main.go         # Entry point
  config/                    # Configuration files
  internal/
    config/                  # Config loader (Viper)
    database/                # Database init
    handler/                 # HTTP handlers
    middleware/              # Auth, CORS, logging, recovery
    model/                   # GORM models
    dto/request/             # Request DTOs
    dto/response/            # Response DTOs
    repository/              # Data access layer
    service/                 # Business logic
    router/                  # Route registration
    scheduler/               # Background tasks
  scripts/                   # External test scripts
  Dockerfile                 # Multi-stage build
  docker-compose.yml         # Server + PostgreSQL
uniapp/                      # UniApp frontend
```

## Features

- User auth (register/login/JWT)
- Product categories & search
- Shopping cart & orders
- Payment integration (mock/wechat/alipay)
- Coupons & flash sales
- Points system
- Merchant applications
- Browse history & favorites
- Address management
- Admin: user management, merchant review, payment refund
