# Engineering Specification

## 1. Architecture Overview

```
                        ┌──────────────┐
                        │   Browser    │
                        │  Vue 3 SPA   │
                        └──────┬───────┘
                               │ HTTPS
                               ▼
                        ┌──────────────┐
                        │   Nginx      │  (Docker, serves frontend + reverse proxy)
                        └──────┬───────┘
                               │ /api/*
                               ▼
                        ┌──────────────┐
                        │   Go API     │  (Chi router, port 8080)
                        │              │
                        │  Middleware:  │
                        │  - CORS      │
                        │  - Auth/JWT  │
                        │  - Logging   │
                        │  - Recovery  │
                        └──────┬───────┘
                               │
                    ┌──────────┼──────────┐
                    ▼          ▼          ▼
             ┌───────────┐ ┌────────┐ ┌────────────┐
             │ PostgreSQL│ │ EBM    │ │ File       │
             │  (sqlc)   │ │ Stub   │ │ Storage    │
             │           │ │ Service│ │ (local vol)│
             └───────────┘ └────────┘ └────────────┘
```

### Key Architectural Decisions

- **Monolith API:** Single Go binary with Chi router. No microservices.
  Modules are separated by Go packages, not by services.
- **sqlc for DB access:** SQL queries are written in `.sql` files; `sqlc`
  generates type-safe Go code. No ORM.
- **EBM is isolated:** The EBM stub is a Go interface + mock implementation
  in its own package. When real RRA docs are available, we swap the
  implementation without touching any other code.
- **No SSR:** The Vue app is a static SPA served by Nginx. All data
  comes from the REST API.

---

## 2. Tech Stack

| Layer | Technology | Version | Why |
|---|---|---|---|
| Frontend framework | Vue 3 (Composition API, `<script setup>`) | 3.x | University requirement |
| Frontend language | TypeScript | 5.x | Type safety |
| Build tool | Vite | 5.x | Fast dev server, good Vue support |
| CSS framework | Tailwind CSS | 3.x | Utility-first, fast styling |
| UI components | PrimeVue | 4.x | Rich component library (DataTable, Dialog, Toast) |
| State management | Pinia | 2.x | Official Vue state management |
| HTTP client | Axios | 1.x | Interceptors for auth token refresh |
| Backend language | Go | 1.22+ | Performance, simplicity |
| HTTP router | chi | 5.x | Lightweight, idiomatic middleware |
| Database | PostgreSQL | 16 | Robust, transactional |
| DB query gen | sqlc | 1.x | Type-safe SQL, no ORM magic |
| DB migrations | golang-migrate | 4.x | Version-controlled schema changes |
| Auth tokens | JWT (golang-jwt) | 5.x | Stateless auth |
| Password hashing | bcrypt | — | Industry standard |
| Containerization | Docker + Docker Compose | — | One-command setup |
| API documentation | Hand-written in `/docs/api/` | — | Keep it simple |

---

## 3. Folder Structure

```
supermarket-management/
│
├── docs/                          # Project documentation
│   ├── product.md
│   ├── ui.md
│   ├── engineering.md
│   └── api/                       # API endpoint documentation per module
│       ├── auth.md
│       ├── products.md
│       └── ...
│
├── frontend/                      # Vue 3 SPA
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── public/
│   │   └── favicon.ico
│   └── src/
│       ├── main.ts                # App entry point
│       ├── App.vue                # Root component
│       ├── router/
│       │   └── index.ts           # Vue Router config with role guards
│       ├── stores/                # Pinia stores
│       │   ├── auth.ts
│       │   ├── products.ts
│       │   ├── cart.ts
│       │   └── ...
│       ├── composables/           # Reusable composition functions
│       │   ├── useAuth.ts
│       │   └── useApi.ts
│       ├── services/              # API call functions (Axios)
│       │   ├── api.ts             # Axios instance + interceptors
│       │   ├── authService.ts
│       │   ├── productService.ts
│       │   └── ...
│       ├── types/                 # TypeScript interfaces
│       │   ├── auth.ts
│       │   ├── product.ts
│       │   └── ...
│       ├── views/                 # Page-level components (one per route)
│       │   ├── LoginView.vue
│       │   ├── DashboardView.vue
│       │   ├── ProductsView.vue
│       │   ├── PosView.vue
│       │   └── ...
│       ├── components/            # Reusable UI components
│       │   ├── layout/
│       │   │   ├── AppSidebar.vue
│       │   │   ├── AppTopbar.vue
│       │   │   └── AppLayout.vue
│       │   ├── products/
│       │   │   ├── ProductTable.vue
│       │   │   └── ProductForm.vue
│       │   ├── pos/
│       │   │   ├── CartPanel.vue
│       │   │   ├── BarcodeInput.vue
│       │   │   └── PayDialog.vue
│       │   └── ...
│       └── assets/                # Static assets (images, fonts)
│
├── backend/                       # Go API server
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # Entry point: load config, connect DB, start server
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go          # Env-based configuration
│   │   ├── database/
│   │   │   ├── db.go              # Connection pool setup
│   │   │   ├── migrations/        # SQL migration files (001_init.up.sql, etc.)
│   │   │   ├── queries/           # sqlc .sql query files, organized by module
│   │   │   │   ├── users.sql
│   │   │   │   ├── products.sql
│   │   │   │   ├── inventory.sql
│   │   │   │   ├── purchases.sql
│   │   │   │   ├── sales.sql
│   │   │   │   └── reports.sql
│   │   │   └── sqlc/              # sqlc generated code (DO NOT EDIT)
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       ├── users.sql.go
│   │   │       └── ...
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT validation + role extraction
│   │   │   ├── logging.go         # Request logging
│   │   │   └── cors.go            # CORS configuration
│   │   ├── auth/
│   │   │   ├── jwt.go             # Token generation, validation, refresh
│   │   │   └── password.go        # bcrypt hash + compare
│   │   ├── handler/               # HTTP handlers (one file per module)
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── category_handler.go
│   │   │   ├── inventory_handler.go
│   │   │   ├── supplier_handler.go
│   │   │   ├── purchase_handler.go
│   │   │   ├── sale_handler.go
│   │   │   ├── receipt_handler.go
│   │   │   ├── dashboard_handler.go
│   │   │   └── report_handler.go
│   │   ├── service/               # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   ├── product_service.go
│   │   │   ├── sale_service.go    # Contains the atomic sale transaction
│   │   │   └── ...
│   │   ├── router/
│   │   │   └── router.go          # Chi route definitions + middleware wiring
│   │   ├── dto/                   # Request/response structs
│   │   │   ├── auth_dto.go
│   │   │   ├── product_dto.go
│   │   │   └── ...
│   │   ├── ebm/                   # EBM stub (isolated)
│   │   │   ├── ebm.go             # Interface definition
│   │   │   └── mock.go            # Mock implementation
│   │   └── validator/
│   │       └── validator.go       # Input validation helpers
│   ├── sqlc.yaml                  # sqlc configuration
│   └── Dockerfile
│
├── docker-compose.yml             # PostgreSQL + Backend + Frontend (Nginx)
├── .env.example                   # Environment variable template
├── Makefile                       # Common commands: migrate, sqlc-generate, test, run
├── seed/                          # Seed data SQL scripts for demo
│   └── seed.sql
└── README.md
```

---

## 4. Database Overview

Only tables needed for Phases 1–6 are defined here. Report queries
(Phase 7) use these same tables.

### Entity-Relationship Diagram

```
┌──────────────┐       ┌───────────────┐       ┌──────────────────┐
│    users     │       │  categories   │       │    products      │
├──────────────┤       ├───────────────┤       ├──────────────────┤
│ id (PK)      │       │ id (PK)       │   ┌──▶│ id (PK)          │
│ name         │       │ name          │   │   │ name             │
│ email (UQ)   │       │ description   │   │   │ barcode (UQ)     │
│ password_hash│       │ created_at    │   │   │ sku              │
│ role         │       │ updated_at    │   │   │ category_id (FK) │──▶ categories
│ is_active    │       └───────────────┘   │   │ unit_price       │
│ created_at   │                           │   │ cost_price       │
│ updated_at   │                           │   │ tax_rate         │
└──────┬───────┘                           │   │ unit_of_measure  │
       │                                   │   │ reorder_level    │
       │ (created_by / cashier_id)         │   │ image_url        │
       │                                   │   │ is_active        │
       ▼                                   │   │ created_at       │
┌──────────────────┐                       │   │ updated_at       │
│    suppliers     │                       │   └──────────────────┘
├──────────────────┤                       │            ▲
│ id (PK)          │                       │            │
│ name             │                       │   ┌────────┴─────────┐
│ contact_person   │                       │   │   inventory      │
│ phone            │                       │   ├──────────────────┤
│ email            │                       │   │ id (PK)          │
│ tin              │                       │   │ product_id (FK)  │
│ address          │                       │   │ quantity         │
│ created_at       │                       │   │ updated_at       │
│ updated_at       │                       │   └──────────────────┘
└──────┬───────────┘                       │
       │                                   │
       ▼                                   │
┌──────────────────┐    ┌──────────────────┤
│ purchase_orders  │    │                  │
├──────────────────┤    │   ┌──────────────┴──────┐
│ id (PK)          │    │   │ stock_movements     │
│ supplier_id (FK) │    │   ├─────────────────────┤
│ status           │    │   │ id (PK)             │
│ total_amount     │    │   │ product_id (FK)     │──▶ products
│ notes            │    │   │ movement_type       │  (purchase_in, sale_out, adjustment)
│ created_by (FK)  │    │   │ quantity            │  (+/−)
│ created_at       │    │   │ reference_id        │  (sale_id or po_id or null)
│ updated_at       │    │   │ reason              │
└──────┬───────────┘    │   │ created_by (FK)     │──▶ users
       │                │   │ created_at          │
       ▼                │   └─────────────────────┘
┌──────────────────┐    │
│ purchase_items   │    │
├──────────────────┤    │
│ id (PK)          │    │
│ po_id (FK)       │──▶ purchase_orders
│ product_id (FK)  │────┘
│ quantity         │
│ unit_cost        │
│ total            │
└──────────────────┘

┌──────────────────┐       ┌──────────────────┐
│     sales        │       │   sale_items     │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │◀──────│ sale_id (FK)     │
│ receipt_number   │       │ id (PK)          │
│ cashier_id (FK)  │──▶usr │ product_id (FK)  │──▶ products
│ subtotal         │       │ product_name     │  (snapshot at time of sale)
│ tax_amount       │       │ quantity         │
│ total_amount     │       │ unit_price       │  (snapshot)
│ payment_method   │       │ tax_rate         │  (snapshot)
│ amount_tendered  │       │ line_total       │
│ change_amount    │       └──────────────────┘
│ status           │  (completed / voided)
│ voided_by (FK)   │  nullable
│ void_reason      │  nullable
│ ebm_signature    │  nullable (mock)
│ created_at       │
└──────────────────┘
```

### Key Design Notes

1. **`inventory` table** holds one row per product with the current
   quantity. It is the denormalized "current stock" for fast reads.
   `stock_movements` is the immutable audit log.

2. **`sale_items` snapshots** `product_name`, `unit_price`, and `tax_rate`
   at the time of sale. This ensures historical receipts remain accurate
   even if the product is later edited or deleted.

3. **`receipt_number`** is a sequential, human-readable number
   (e.g., `RCP-20260929-0001`). Generated by the backend, not the DB
   sequence directly.

4. **`stock_movements.reference_id`** is a polymorphic FK: it points to
   `sales.id` for sale_out, `purchase_orders.id` for purchase_in, or
   is NULL for manual adjustments.

5. **Role is stored as a string enum** in the `users` table:
   `admin`, `manager`, `accountant`, `storekeeper`, `cashier`.

---

## 5. API Conventions

### Base URL

```
/api/v1
```

### Request / Response Format

- All requests and responses use `Content-Type: application/json`.
- Dates are ISO 8601: `2026-09-29T14:30:00Z`.
- Money values are integers in the smallest currency unit (RWF has no
  decimals, so 700 RWF = `700`).

### Standard Response Envelope

**Success (single resource):**
```json
{
  "data": { ... },
  "message": "Product created successfully"
}
```

**Success (list):**
```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 87,
    "total_pages": 5
  }
}
```

**Error:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Barcode is required",
    "details": [
      { "field": "barcode", "message": "must not be empty" }
    ]
  }
}
```

### HTTP Methods & Status Codes

| Action | Method | Success Code |
|---|---|---|
| List / Get | GET | 200 |
| Create | POST | 201 |
| Update | PUT or PATCH | 200 |
| Delete | DELETE | 204 |
| Login | POST | 200 |

| Error | Code |
|---|---|
| Validation error | 400 |
| Bad credentials | 401 |
| Forbidden (wrong role) | 403 |
| Not found | 404 |
| Conflict (duplicate barcode) | 409 |
| Server error | 500 |

### Pagination

Query params: `?page=1&per_page=20`
Default: page 1, 20 items per page. Max per_page: 100.

### Filtering & Sorting

Query params: `?search=milk&category_id=3&status=active&sort=name&order=asc`

### Endpoint Naming (examples)

```
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh

GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}

GET    /api/v1/categories
POST   /api/v1/categories
PUT    /api/v1/categories/{id}
DELETE /api/v1/categories/{id}

GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/{id}
PUT    /api/v1/products/{id}
GET    /api/v1/products/barcode/{barcode}

GET    /api/v1/inventory
POST   /api/v1/inventory/adjustments
GET    /api/v1/inventory/movements

GET    /api/v1/suppliers
POST   /api/v1/suppliers
PUT    /api/v1/suppliers/{id}

GET    /api/v1/purchases
POST   /api/v1/purchases
GET    /api/v1/purchases/{id}
PATCH  /api/v1/purchases/{id}/status

POST   /api/v1/sales
GET    /api/v1/sales
GET    /api/v1/sales/{id}
POST   /api/v1/sales/{id}/void

GET    /api/v1/receipts/{sale_id}

GET    /api/v1/dashboard
GET    /api/v1/reports/sales
GET    /api/v1/reports/inventory
GET    /api/v1/reports/purchases
```

---

## 6. Authentication & Authorization

### Auth Flow

```
┌────────┐                          ┌──────────┐                    ┌────────┐
│ Client │                          │  Go API  │                    │   DB   │
└───┬────┘                          └────┬─────┘                    └───┬────┘
    │  POST /auth/login                  │                              │
    │  { email, password }               │                              │
    │───────────────────────────────────▶│                              │
    │                                    │  SELECT user by email        │
    │                                    │─────────────────────────────▶│
    │                                    │  user row                    │
    │                                    │◀─────────────────────────────│
    │                                    │                              │
    │                                    │  bcrypt.Compare(password,    │
    │                                    │    password_hash)            │
    │                                    │                              │
    │  200 {                             │                              │
    │    access_token (15min),           │                              │
    │    refresh_token (7d)              │                              │
    │  }                                 │                              │
    │◀──────────────────────────────────│                              │
    │                                    │                              │
    │  GET /products                     │                              │
    │  Authorization: Bearer <token>     │                              │
    │───────────────────────────────────▶│                              │
    │                                    │  middleware: validate JWT,   │
    │                                    │  extract role from claims    │
    │                                    │  check role ∈ allowed_roles  │
    │                                    │                              │
    │  200 { data: [...] }               │                              │
    │◀──────────────────────────────────│                              │
```

### JWT Claims

```json
{
  "sub": "user-uuid",
  "email": "marie@store.rw",
  "role": "cashier",
  "exp": 1727625000,
  "iat": 1727624100
}
```

### Middleware Stack

1. **CORS** — allow frontend origin.
2. **Request Logger** — method, path, status, duration.
3. **Recovery** — catch panics, return 500.
4. **Auth (per-route group):**
   - Public routes: `/auth/login`, `/auth/refresh` — no middleware.
   - Protected routes: `RequireAuth()` validates JWT, injects user
     into context.
   - Role-restricted routes: `RequireRole("admin", "manager")` checks
     the role claim.

### Token Storage (Frontend)

- **Access token:** Stored in a Pinia store (memory only). Lost on
  page refresh — the app calls `/auth/refresh` on startup.
- **Refresh token:** Sent as an `httpOnly`, `Secure`, `SameSite=Strict`
  cookie. Not accessible to JavaScript.

---

## 7. Sale Transaction Design

This is the most critical piece of backend logic. It must be **atomic**:
either the entire sale succeeds or nothing changes.

### Sequence

```
POST /api/v1/sales
{
  "items": [
    { "product_id": "uuid", "quantity": 2 },
    { "product_id": "uuid", "quantity": 1 }
  ],
  "payment_method": "cash",
  "amount_tendered": 5000
}
```

### Backend Steps (inside a single DB transaction)

```go
// Pseudocode — actual implementation uses sqlc-generated functions

func (s *SaleService) CreateSale(ctx, req) (*SaleResponse, error) {

    tx, _ := db.BeginTx(ctx, nil)
    defer tx.Rollback()

    qtx := s.queries.WithTx(tx)

    // 1. Generate receipt number
    receiptNum := generateReceiptNumber(ctx, qtx)  // e.g. RCP-20260929-0042

    // 2. Look up each product, calculate totals
    var subtotal, taxTotal int64
    var saleItems []SaleItem

    for _, item := range req.Items {
        product, err := qtx.GetProduct(ctx, item.ProductID)
        // → 404 if not found, 400 if inactive

        // 3. Check stock
        inv, err := qtx.GetInventory(ctx, item.ProductID)
        if inv.Quantity < item.Quantity {
            return error "Insufficient stock for {product.Name}"
        }

        lineTotal := product.UnitPrice * item.Quantity
        lineTax   := lineTotal * product.TaxRate / 100

        saleItems = append(saleItems, SaleItem{
            ProductID:   product.ID,
            ProductName: product.Name,      // snapshot
            Quantity:    item.Quantity,
            UnitPrice:   product.UnitPrice,  // snapshot
            TaxRate:     product.TaxRate,    // snapshot
            LineTotal:   lineTotal,
        })

        subtotal += lineTotal
        taxTotal += lineTax

        // 4. Decrement inventory
        qtx.DecrementInventory(ctx, product.ID, item.Quantity)

        // 5. Record stock movement
        qtx.CreateStockMovement(ctx, StockMovement{
            ProductID:    product.ID,
            MovementType: "sale_out",
            Quantity:     -item.Quantity,
            ReferenceID:  saleID,
            CreatedBy:    currentUser.ID,
        })
    }

    totalAmount := subtotal + taxTotal
    changeAmount := req.AmountTendered - totalAmount

    // 6. Create sale record
    sale := qtx.CreateSale(ctx, Sale{
        ReceiptNumber: receiptNum,
        CashierID:     currentUser.ID,
        Subtotal:      subtotal,
        TaxAmount:     taxTotal,
        TotalAmount:   totalAmount,
        PaymentMethod: req.PaymentMethod,
        AmountTendered: req.AmountTendered,
        ChangeAmount:  changeAmount,
        Status:        "completed",
    })

    // 7. Insert sale items
    for _, item := range saleItems {
        qtx.CreateSaleItem(ctx, sale.ID, item)
    }

    // 8. Call EBM stub (outside tx if it's external; inside if mock)
    ebmSig := s.ebmService.ReportSale(sale)
    qtx.UpdateSaleEBM(ctx, sale.ID, ebmSig)

    // 9. Commit
    tx.Commit()

    return &SaleResponse{Sale: sale, Items: saleItems, EBMSignature: ebmSig}
}
```

### Why This Design

- **Single transaction:** If any step fails (insufficient stock, DB error),
  the entire sale rolls back. No partial inventory decrements.
- **Snapshots in `sale_items`:** Product name and price are copied at
  sale time. Even if the product is later renamed or repriced, old
  receipts remain correct.
- **Stock movement log:** Every inventory change is recorded with who,
  when, why, and a reference to the originating sale or PO.
- **EBM is pluggable:** `s.ebmService` is an interface. The mock returns
  a fake signature. A future real implementation swaps in without
  changing the sale flow.

---

## 8. Testing Strategy

Testing is organized into four tiers. Each module (auth, products,
inventory, purchases, sales, reports) follows the same pattern.

### Tier 1 — Unit Tests (Go)

**Scope:** Pure business logic with no database or HTTP dependencies.
**Framework:** Go standard `testing` package + `testify/assert`.
**Location:** `*_test.go` files alongside the code they test.

| Module | What to Test |
|---|---|
| Auth | Password hashing/comparing, JWT generation/validation/expiry, role extraction from claims |
| Products | Price calculation helpers, barcode format validation, tax computation |
| Inventory | Stock sufficiency check logic, reorder-level alert logic |
| Sales | Receipt number generation, change calculation, line total / tax calculation |
| EBM stub | Mock returns expected signature format |
| Validators | Input validation functions: required fields, email format, number ranges |

**Run with:**
```bash
go test ./internal/... -short -count=1
```

### Tier 2 — Database / Integration Tests (Go + PostgreSQL)

**Scope:** sqlc-generated queries against a real test database.
**Setup:** Docker Compose starts a test PostgreSQL container. Migrations
run before tests. Each test runs in a transaction that is rolled back
(or uses a per-test schema).

| Module | What to Test |
|---|---|
| Users | Create user, duplicate email rejected, find by email, list with pagination |
| Categories | CRUD, delete blocked if products exist |
| Products | CRUD, unique barcode constraint, filter by category, barcode lookup |
| Inventory | Decrement quantity, prevent negative stock (CHECK constraint), stock movement insert |
| Purchases | PO creation with items, status transitions, goods receiving increments stock |
| Sales | Full sale transaction: sale + items + inventory decrement + movement log in one tx; rollback on insufficient stock |

**Run with:**
```bash
# Requires test DB running (docker compose up db-test)
go test ./internal/database/... -count=1 -tags=integration
```

### Tier 3 — API / HTTP Tests (Go)

**Scope:** Full HTTP request → response cycle, including middleware (auth,
role checks, validation).
**Framework:** Go `net/http/httptest` + `chi` test server.
**Setup:** Use the real router with a test database (or mock service layer
for pure handler tests).

| Module | What to Test |
|---|---|
| Auth | Login success/failure, token in response, refresh flow, expired token → 401 |
| Users | Admin can create user, cashier cannot (403), duplicate email → 409 |
| Products | Manager can CRUD, storekeeper can only GET, unauthenticated → 401 |
| Inventory | Storekeeper can adjust, cashier cannot, adjustment creates movement log |
| Purchases | Create PO, receive goods → stock updated, cancel PO |
| Sales | Create sale → 201 + receipt data, insufficient stock → 400, void sale → inventory restored |
| Receipts | GET receipt returns correct items and totals |
| Dashboard | Returns today's stats, respects role access |

**What each test checks:**
- HTTP status code
- Response body structure (correct fields, types)
- Role-based access (403 for unauthorized roles)
- Validation errors (400 with field-level messages)
- Side effects (e.g., inventory changed after sale)

**Run with:**
```bash
go test ./internal/handler/... -count=1 -tags=integration
```

### Tier 4 — Manual Testing Checklist (per module)

These are performed by the developer during each phase, before moving
to the next phase. Keep a checklist in the PR or in a testing log.

#### Phase 1 — Auth & Users

- [ ] Login with valid credentials → redirected to correct page per role
- [ ] Login with wrong password → error message shown
- [ ] Login with disabled account → error message shown
- [ ] Access /products as unauthenticated → redirected to login
- [ ] Admin creates a new user → user appears in list
- [ ] Admin disables a user → user cannot log in
- [ ] Non-admin tries /users page → 403 or redirect

#### Phase 2 — Products & Categories

- [ ] Create category → appears in list
- [ ] Create product with barcode → appears in list
- [ ] Duplicate barcode → error shown
- [ ] Search product by name → filtered results
- [ ] Filter by category → correct products shown
- [ ] Edit product price → price updated
- [ ] Deactivate product → hidden from POS (but visible in product list)

#### Phase 3 — Inventory

- [ ] After product creation, inventory row exists with qty = 0
- [ ] Manual stock adjustment → qty changes, movement logged
- [ ] Low-stock product appears in alerts
- [ ] Stock movement log shows all adjustments with user and reason

#### Phase 4 — Purchases

- [ ] Create supplier → appears in list
- [ ] Create PO with 3 line items → saved as Draft
- [ ] Mark PO as Received → inventory for each product increases
- [ ] Stock movements logged as "purchase_in"
- [ ] Cancel a PO → status changes, no stock impact

#### Phase 5 — POS & Sales

- [ ] Cashier logs in → POS screen opens, barcode input focused
- [ ] Scan valid barcode → product appears in cart
- [ ] Scan same barcode again → quantity increments
- [ ] Scan invalid barcode → toast error, input cleared
- [ ] Change quantity in cart → totals update
- [ ] Remove item from cart → totals update
- [ ] Hold sale → cart saved, new empty cart
- [ ] Recall held sale → cart restored
- [ ] Complete sale with cash → receipt shown, cart cleared
- [ ] Check inventory → quantities decreased correctly
- [ ] Attempt sale with insufficient stock → error message

#### Phase 6 — Receipts & EBM

- [ ] Receipt after sale shows all required fields (store, TIN, items, tax, total)
- [ ] Receipt has mock EBM QR code / signature
- [ ] Print receipt → browser print dialog opens
- [ ] Manager can search and reprint a past receipt
- [ ] Void sale → status changes, inventory restored, void reason recorded

#### Phase 7 — Dashboard & Reports

- [ ] Dashboard shows today's revenue (matches sum of today's sales)
- [ ] Top 5 products reflect today's sales
- [ ] Low-stock alerts show correct products
- [ ] Sales report with date filter → correct totals
- [ ] Export CSV → file downloads, data matches screen
- [ ] Inventory report → shows current stock with value at cost

#### Phase 8 — End-to-End Demo

- [ ] Run through the complete Demo Story (product.md §6) without errors
- [ ] `docker compose up` from scratch → system works
- [ ] Seed data loaded correctly
- [ ] All five roles can log in and see only their permitted pages
- [ ] Demo completes in under 10 minutes

---

## 9. Phase-by-Phase Build Order

| Phase | Backend | Frontend | Database |
|---|---|---|---|
| **1. Scaffold & Auth** | Go project, Chi router, config, CORS, JWT middleware, auth handler, user handler | Vue + Vite + Tailwind + PrimeVue scaffold, Login page, User management page, route guards | `users` table, migrations setup |
| **2. Products & Categories** | Category CRUD handler, Product CRUD handler, barcode lookup | Products page, Categories page, product form | `categories`, `products` tables |
| **3. Inventory** | Inventory endpoints, stock adjustment, movement log | Inventory page (tabs: levels, movements, adjustments) | `inventory`, `stock_movements` tables |
| **4. Purchases** | Supplier CRUD, PO CRUD, receive goods (stock-in) | Suppliers page, Purchases page, PO detail | `suppliers`, `purchase_orders`, `purchase_items` tables |
| **5. POS & Sales** | Sale creation (atomic tx), void sale | POS full-screen, cart, pay dialog | `sales`, `sale_items` tables |
| **6. Receipts & EBM** | Receipt generation, EBM stub interface + mock | Receipt view, print/PDF | EBM signature column on `sales` |
| **7. Dashboard & Reports** | Dashboard aggregation queries, report endpoints, CSV export | Dashboard page, Reports page | No new tables — queries on existing |
| **8. Polish & Deploy** | Seed script, Dockerfile, docker-compose | UI polish, loading states, error handling | Seed data SQL |

