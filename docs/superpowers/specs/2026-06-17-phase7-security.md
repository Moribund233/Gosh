# Phase 7.2 — Security & Input Protection

**Scope:** Cart quantity cap, upload file size limit, XSS protection.

## Items

### 1. Cart quantity max validation

**Current:** `AddCartRequest.Quantity` has `binding:"required,min=1"` — no upper bound.

**Change:** Add `max=99` to the binding tag. `CartMaxQuantity = 99` constant already defined.

| File | Change |
|------|--------|
| `internal/dto/request/cart.go:5` | `Quantity int \`json:"quantity" binding:"required,min=1,max=99"\`` |

### 2. Upload MaxMultipartMemory

**Current:** `gin.Engine` has no `MaxMultipartMemory` set (default 32MB). Upload handler checks size after receiving the file, but oversized request bodies can consume server resources before reaching the check.

**Change:** Set `r.MaxMultipartMemory = 10 << 20` in `router.go:New()` — rejects multipart requests > 10MB at the Gin framework level.

| File | Change |
|------|--------|
| `internal/router/router.go:34` | Add `r.MaxMultipartMemory = 10 << 20` |

### 3. Pagination defaults

**Status:** ✅ Already done across all endpoints. Skipping.

### 4. XSS output/input encoding

**Approach:** Gin middleware + utility package.

**Middleware (primary defense):** `internal/middleware/sanitize.go` — intercepts POST/PUT/PATCH requests, reads JSON body, recursively HTML-escapes all string fields, replaces body for downstream handlers.

**Utility (targeted defense):** `pkg/sanitize/sanitize.go` — `String(s string) string` that calls `html.EscapeString`, usable in service layer for edge cases.

| File | Change |
|------|--------|
| `pkg/sanitize/sanitize.go` | Create: `String(s string) string` |
| `internal/middleware/sanitize.go` | Create: Gin middleware applying sanitize to JSON body |
| `internal/router/router.go` | Register middleware |

## Testing

- Cart: existing test `TestAddCart_Success` should still pass with `max=99` (test uses quantity 1)
- XSS middleware: unit test verifying `<script>` tags are escaped in request body
- Upload: no new test needed (Gin-level config)

## Files changed

- `server/internal/dto/request/cart.go` — 1 line
- `server/internal/router/router.go` — 2 lines
- `server/pkg/sanitize/sanitize.go` — new, ~10 lines
- `server/internal/middleware/sanitize.go` — new, ~50 lines
