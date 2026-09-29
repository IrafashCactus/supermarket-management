# UI Specification

## 1. Global Layout

```
┌──────────────────────────────────────────────────┐
│  Top Bar: Logo / App Name    [User ▼] [Logout]   │
├────────┬─────────────────────────────────────────┤
│        │                                         │
│  Side  │           Main Content Area             │
│  Nav   │                                         │
│        │                                         │
│        │                                         │
├────────┴─────────────────────────────────────────┤
│  Status Bar (optional): connection, version       │
└──────────────────────────────────────────────────┘
```

- **Sidebar** collapses on small screens (hamburger menu).
- **POS screen** is an exception: it uses a **full-screen layout** with no sidebar, optimized for speed.
- Navigation items are filtered by the logged-in user's role.

### Sidebar Navigation by Role

| Menu Item | Admin | Manager | Accountant | Storekeeper | Cashier |
|---|:---:|:---:|:---:|:---:|:---:|
| Dashboard | ✓ | ✓ | ✓ | — | — |
| Products | ✓ | ✓ | — | view | — |
| Categories | ✓ | ✓ | — | — | — |
| Inventory | ✓ | ✓ | — | ✓ | — |
| Suppliers | ✓ | ✓ | — | ✓ | — |
| Purchases | ✓ | ✓ | — | ✓ | — |
| POS | — | — | — | — | ✓ |
| Sales | ✓ | ✓ | ✓ | — | own |
| Reports | ✓ | ✓ | ✓ | — | — |
| Users | ✓ | — | — | — | — |

---

## 2. Screen-by-Screen Specification

---

### 2.1 Login Screen

**Route:** `/login`
**Access:** Public (unauthenticated)

| Element | Details |
|---|---|
| Logo & app name | Centered, above form |
| Email field | Text input, required |
| Password field | Password input, required, toggle visibility |
| "Sign In" button | Submits credentials, shows spinner while loading |
| Error message | Inline alert for invalid credentials |

**Behavior:**
- On success → redirect to the user's default landing page (Cashier → POS, others → Dashboard).
- JWT stored in memory (not localStorage). Refresh token in httpOnly cookie.

---

### 2.2 Dashboard

**Route:** `/dashboard`
**Access:** Admin, Manager, Accountant

| Section | Content |
|---|---|
| KPI cards (top row) | Today's Revenue, Today's Sales Count, Low-Stock Items Count, Products Count |
| Revenue chart | Bar or line chart — sales revenue for the last 7 days |
| Top 5 products | Small table: product name, units sold today |
| Low-stock alerts | List of products below reorder level with current qty and reorder level |
| Recent sales | Table of last 10 sales: receipt #, time, cashier, total |

---

### 2.3 Products Screen

**Route:** `/products`
**Access:** Admin, Manager (full CRUD); Storekeeper (view only)

| Element | Details |
|---|---|
| Search bar | Search by name or barcode, instant filter |
| Category filter | Dropdown to filter by category |
| Status filter | Active / Inactive / All |
| Product table | Columns: Image thumbnail, Name, Barcode, Category, Price, Cost, Stock Qty, Status |
| "Add Product" button | Opens create dialog/drawer |
| Row actions | Edit, Toggle active/inactive |

**Add / Edit Product Dialog:**

| Field | Type | Required |
|---|---|:---:|
| Name | Text | ✓ |
| Barcode | Text | ✓ |
| SKU | Text | — |
| Category | Dropdown | ✓ |
| Unit Price (selling) | Number | ✓ |
| Cost Price | Number | ✓ |
| Tax Rate (%) | Number (default 18) | ✓ |
| Unit of Measure | Dropdown (piece, kg, litre, pack) | ✓ |
| Reorder Level | Number | — |
| Image | File upload | — |
| Active | Toggle | ✓ |

---

### 2.4 Categories Screen

**Route:** `/categories`
**Access:** Admin, Manager

| Element | Details |
|---|---|
| Category table | Columns: Name, Description, Product Count |
| "Add Category" button | Opens inline form or dialog |
| Row actions | Edit, Delete (only if no products linked) |

---

### 2.5 Inventory Screen

**Route:** `/inventory`
**Access:** Admin, Manager, Storekeeper

| Tab | Content |
|---|---|
| **Stock Levels** | Table: Product, Barcode, Category, Current Qty, Reorder Level, Status (OK / Low / Out). Sortable, filterable. |
| **Stock Movements** | Log table: Date/Time, Product, Movement Type (Purchase In, Sale Out, Adjustment), Quantity (+/−), User, Reason. Filterable by date range and type. |
| **Adjustments** | "New Adjustment" button → dialog: select product, +/− quantity, reason (Damage, Theft, Correction, Other), notes. Requires manager approval if configured. |
| **Expiry Alerts** | Table of batches expiring within 30 days: Product, Batch, Expiry Date, Qty remaining. |

---

### 2.6 Suppliers Screen

**Route:** `/suppliers`
**Access:** Admin, Manager, Storekeeper

| Element | Details |
|---|---|
| Supplier table | Columns: Name, Contact Person, Phone, Email, TIN |
| "Add Supplier" button | Opens dialog |
| Row actions | Edit, View purchases |

---

### 2.7 Purchases Screen

**Route:** `/purchases`
**Access:** Admin, Manager, Storekeeper

| Element | Details |
|---|---|
| Purchase orders table | Columns: PO #, Supplier, Date, Status (Draft / Sent / Received / Cancelled), Total, Created By |
| Status filter | Dropdown |
| "New Purchase Order" button | Opens PO creation page |

**PO Creation / Detail Page:**

| Section | Content |
|---|---|
| Header | Supplier (dropdown), Date, Notes |
| Line items table | Product (searchable dropdown), Quantity, Unit Cost, Line Total. Add/remove rows. |
| Totals | Subtotal, Tax, Grand Total |
| Actions | Save Draft, Mark as Sent, Receive (storekeeper confirms quantities → stock-in), Cancel |

---

### 2.8 POS Screen ⭐

**Route:** `/pos`
**Access:** Cashier (full screen, no sidebar)

```
┌─────────────────────────────────────────────────────────────────┐
│  [≡ Menu]   SUPERMARKET POS            Cashier: Marie   [Exit] │
├─────────────────────────────────┬───────────────────────────────┤
│                                 │                               │
│   CART                          │   PRODUCT INPUT               │
│                                 │                               │
│   ┌───────────────────────────┐ │   ┌─────────────────────────┐ │
│   │ # │ Product  │Qty│ Price │ │   │ 🔍 Scan or search...    │ │
│   │───│──────────│───│───────│ │   └─────────────────────────┘ │
│   │ 1 │ Milk 500 │ 1 │  700  │ │                               │
│   │ 2 │ Bread    │ 2 │ 1600  │ │   Recent / Quick-add grid     │
│   │ 3 │ Fanta 1L │ 1 │  800  │ │   (optional, phase 8 polish)  │
│   │   │          │   │       │ │                               │
│   └───────────────────────────┘ │                               │
│                                 │                               │
│   Subtotal:          3,100 RWF  │                               │
│   VAT (18%):           558 RWF  │                               │
│   ─────────────────────────     │                               │
│   TOTAL:             3,658 RWF  │                               │
│                                 │                               │
│   [Hold Sale]  [Clear]  [Pay ▶] │                               │
│                                 │                               │
├─────────────────────────────────┴───────────────────────────────┤
│  Held Sales: (2)  │  Shift Total: 45,200 RWF  │  Sales: 12     │
└─────────────────────────────────────────────────────────────────┘
```

**Pay Dialog (modal):**

| Element | Details |
|---|---|
| Total due | Large, prominent |
| Payment method | Radio buttons: Cash, Card, Mobile Money |
| Amount tendered | Number input (pre-filled with total) |
| Change due | Auto-calculated |
| "Complete Sale" button | Submits the sale |
| "Cancel" | Returns to cart |

**After sale completes:**
- Receipt overlay appears (see 2.9).
- Cart clears, focus returns to barcode input.

---

### 2.9 Receipt View

**Route:** `/receipts/:id` (standalone) — also shown as overlay after POS sale

| Section | Content |
|---|---|
| Header | Store name, address, TIN, phone |
| Receipt info | Receipt #, Date/Time, Cashier name |
| Items table | Product name, Qty, Unit Price, Line Total |
| Totals | Subtotal, VAT amount, Grand Total |
| Payment | Method, Amount Tendered, Change |
| Footer | "Thank you" message, mock EBM QR code / signature |
| Actions | Print (browser print), Download PDF, Close |

---

### 2.10 Sales Screen

**Route:** `/sales`
**Access:** Admin, Manager, Accountant (all sales); Cashier (own sales only)

| Element | Details |
|---|---|
| Sales table | Columns: Receipt #, Date/Time, Cashier, Items Count, Total, Payment Method, Status (Completed / Voided) |
| Filters | Date range, Cashier (dropdown), Payment method, Status |
| Row click | Opens receipt view |
| Void action | Manager only; prompts for reason; reverses inventory |

---

### 2.11 Reports Screen

**Route:** `/reports`
**Access:** Admin, Manager, Accountant

| Report | Content |
|---|---|
| **Sales Report** | Date range picker. Table of sales with totals. Summary row: total revenue, total tax, sale count. Export CSV. |
| **Inventory Report** | Current stock with value at cost. Filter by category, stock status. Export CSV. |
| **Purchase Report** | Spending by supplier and date range. Export CSV. |
| **Shift Summary** | Select cashier + date. Shows sales count, total revenue, payment method breakdown. |

---

### 2.12 Users Screen

**Route:** `/users`
**Access:** Admin only

| Element | Details |
|---|---|
| Users table | Columns: Name, Email, Role, Status (Active/Disabled), Last Login |
| "Add User" button | Opens dialog: name, email, password, role, active toggle |
| Row actions | Edit, Reset password, Disable/Enable |

---

## 3. POS Cashier Flow (Step by Step)

This is the critical path — the "wow moment" of the demo.

```
┌─────────────┐
│ Cashier logs │
│ in           │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ POS screen  │  ← full-screen, barcode input auto-focused
│ opens       │
└──────┬──────┘
       │
       ▼
┌──────────────────┐     ┌──────────────────┐
│ Scan / type      │────▶│ Product found?   │
│ barcode          │     └────┬────────┬────┘
└──────────────────┘          │ Yes    │ No
                              ▼        ▼
                     ┌────────────┐  ┌──────────────┐
                     │ Add to cart│  │ Show "Product │
                     │ (qty +1)  │  │ not found"   │
                     └────┬───────┘  │ error toast  │
                          │          └──────────────┘
                          ▼
                   ┌──────────────┐
                   │ Cart updated │  ← subtotal, tax, total recalculated
                   │              │
                   └──────┬───────┘
                          │
              ┌───────────┼────────────┐
              ▼           ▼            ▼
     ┌──────────┐  ┌───────────┐  ┌──────────┐
     │ Scan     │  │ Modify    │  │ Click    │
     │ another  │  │ qty /     │  │ [Pay ▶]  │
     │ product  │  │ remove    │  │          │
     └──────────┘  └───────────┘  └────┬─────┘
                                       │
                                       ▼
                                ┌──────────────┐
                                │ Pay dialog   │
                                │ opens        │
                                └──────┬───────┘
                                       │
                                       ▼
                                ┌──────────────┐
                                │ Select method│
                                │ Enter amount │
                                │ tendered     │
                                └──────┬───────┘
                                       │
                                       ▼
                                ┌──────────────────────────────┐
                                │ Click "Complete Sale"        │
                                │                              │
                                │ Backend (atomic transaction): │
                                │  1. Create sale + line items  │
                                │  2. Decrement inventory      │
                                │  3. Generate receipt number   │
                                │  4. Call EBM stub            │
                                │  5. Return receipt data      │
                                └──────────────┬───────────────┘
                                               │
                                               ▼
                                        ┌──────────────┐
                                        │ Receipt      │
                                        │ overlay      │
                                        │ displayed    │
                                        └──────┬───────┘
                                               │
                                      ┌────────┴────────┐
                                      ▼                 ▼
                               ┌───────────┐     ┌───────────┐
                               │ Print /   │     │ Close     │
                               │ Download  │     │           │
                               └───────────┘     └─────┬─────┘
                                                       │
                                                       ▼
                                                ┌──────────────┐
                                                │ Cart cleared │
                                                │ Focus back   │
                                                │ to barcode   │
                                                │ input        │
                                                └──────────────┘
```

### Key UX Details

1. **Auto-focus:** When the POS screen loads, the barcode input is immediately focused. After every action (add, pay, clear), focus returns to the barcode input.

2. **Keyboard-driven:** The entire POS flow should work without a mouse. Barcode scanners behave like keyboards (type digits + Enter).

3. **Speed:** No unnecessary modals or confirmations for adding items. Scan → appears in cart instantly.

4. **Quantity handling:** Scanning the same barcode again increments the quantity (+1). The cashier can also manually edit the quantity in the cart.

5. **Error handling:** If a barcode is not found, show a brief toast (not a blocking modal) and re-focus the input.

6. **Hold / Recall:** Cashier can park the current cart (e.g., customer forgot wallet), serve another customer, then recall the held sale.

7. **Offline resilience (stretch):** If the network drops mid-scan, the cart is local. Sale submission requires connectivity.

