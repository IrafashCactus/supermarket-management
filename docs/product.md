# Product Specification

## 1. Purpose

A Supermarket Management & POS System designed for a single-branch
supermarket in Rwanda. The system covers the full lifecycle of goods —
from supplier purchase to customer sale — and produces fiscal receipts
that comply with Rwanda Revenue Authority (RRA) formatting requirements.

This is a university capstone project. The goal is a working, deployable
system that can be demonstrated end-to-end: a cashier scans a product,
completes a sale, prints a receipt, and inventory is updated instantly.

---

## 2. Users & Roles

| Role | Description | Key Permissions |
|---|---|---|
| **Administrator** | System owner. Manages users, roles, and global settings. | Create/edit/disable users, assign roles, view audit logs, full system access |
| **Manager** | Store manager. Oversees products, pricing, and performance. | CRUD products & categories, set prices, approve stock adjustments, view dashboard & reports |
| **Accountant** | Financial oversight. Reviews sales data and tax invoices. | View sales reports, export financial summaries, view/reprint receipts, view dashboard |
| **Storekeeper** | Warehouse / stock room operator. | Receive goods from suppliers, record stock-in/stock-out, run stock counts, track expiry dates |
| **Cashier** | Point-of-sale operator. | Open POS session, scan/search products, process sales, generate receipts, view own shift summary |

### Permission Principles

- Every role has an implicit "view own profile / change own password" permission.
- Roles are enforced server-side via middleware on every API request.
- The frontend hides UI elements the user cannot access, but the backend is the source of truth.
- An Administrator can do everything any other role can do.

---

## 3. Features per Module

### Phase 1 — Project Scaffold & Authentication

| Feature | Details |
|---|---|
| User login | Email + password, returns JWT access token + refresh token |
| User management (Admin) | CRUD users, assign one role per user, disable accounts |
| Password reset | Admin-initiated password reset (no email flow in v1) |
| Session security | JWT with expiry, refresh token rotation, role claim in token |

### Phase 2 — Products & Categories

| Feature | Details |
|---|---|
| Category management | CRUD categories (name, description); products belong to one category |
| Product management | CRUD products: name, barcode, SKU, category, unit price, cost price, tax rate, unit of measure, image (optional), active/inactive flag |
| Barcode lookup | Search product by barcode (used by POS and inventory) |
| Product search | Search/filter by name, category, barcode, status |
| Bulk import (stretch) | CSV upload to seed product catalog |

### Phase 3 — Inventory Management

| Feature | Details |
|---|---|
| Stock tracking | Each product has a current quantity; quantities update on purchase-in and sale-out |
| Stock adjustments | Manual adjust (damage, theft, correction) with a reason code and who approved |
| Low-stock alerts | Configurable reorder-level per product; dashboard shows products below threshold |
| Expiry tracking | Optional expiry date on stock batches; alerts for items expiring within N days |
| Stock movement log | Immutable log of every quantity change: type (purchase, sale, adjustment), quantity, timestamp, user |

### Phase 4 — Purchases & Suppliers

| Feature | Details |
|---|---|
| Supplier management | CRUD suppliers: name, contact info, TIN |
| Purchase orders | Create PO with line items (product, quantity, unit cost); status workflow: Draft → Sent → Received → Cancelled |
| Goods receiving | Storekeeper confirms received items against PO; stock-in is recorded, inventory incremented |
| Purchase history | Filter/search past purchases by supplier, date range, status |

### Phase 5 — POS & Sales

| Feature | Details |
|---|---|
| POS interface | Full-screen cashier view; scan barcode or search product; cart with line items |
| Cart operations | Add item, change quantity, remove item; real-time subtotal, tax, and total |
| Payment recording | Record payment method (Cash, Card, Mobile Money — labels only, no gateway) and amount tendered; calculate change |
| Sale completion | Atomic transaction: create sale record + decrement inventory for every line item |
| Hold / recall sale | Cashier can park a sale and recall it later (single cashier per terminal) |
| Void / cancel sale | Manager-authorized void of a completed sale; inventory is restored |

### Phase 6 — Receipts & EBM Stub

| Feature | Details |
|---|---|
| Receipt generation | Generate receipt with: store name, TIN, date/time, items, totals, tax breakdown, receipt number |
| Receipt PDF / print | Render receipt as thermal-printer-friendly HTML or downloadable PDF |
| EBM stub service | Isolated service interface that simulates sending invoice data to RRA and returns a mock EBM signature/QR code. Real integration deferred to when official API docs are available. |
| Receipt reprint | Accountant or Manager can look up and reprint any receipt by sale ID |

### Phase 7 — Dashboard & Reports

| Feature | Details |
|---|---|
| Dashboard | Today's revenue, number of sales, top 5 products, low-stock alerts, recent sales list |
| Sales report | Filter by date range, cashier, payment method; totals and line-item detail; export CSV |
| Inventory report | Current stock levels, value at cost, items below reorder level |
| Purchase report | Spending by supplier, by date range |
| Shift summary | Per-cashier sales total for a given shift/day |

### Phase 8 — Polish & Deployment

| Feature | Details |
|---|---|
| UI polish | Responsive layout, loading states, error handling, toast notifications |
| Dockerized deployment | Single `docker compose up` starts frontend, backend, and database |
| Seed data | Script to populate demo categories, products, suppliers, and a test user per role |
| Final testing | End-to-end walkthrough of the demo story |

---

## 4. MVP Scope (What We Ship)

The MVP is **Phases 1–7 fully working**. Phase 8 is polish.
The absolute minimum for a successful demo:

1. Admin can log in and create a cashier user.
2. Manager can add categories and products with barcodes.
3. Storekeeper can receive a purchase and stock goes up.
4. Cashier can scan a barcode, complete a sale, and print a receipt.
5. Inventory decreases after the sale.
6. Dashboard shows today's revenue and the sale that just happened.

---

## 5. Out of Scope (v1)

| Feature | Reason |
|---|---|
| Loyalty / rewards points | Deferred to v2 |
| Real payment gateway (bank, MTN MoMo, Airtel Money) | Deferred; we record the method label only |
| Multi-branch / multi-store | Single-store system for now |
| Camera-based barcode scanning | Requires mobile/WebRTC work; we use keyboard-wedge scanners or manual entry |
| Real EBM / RRA API integration | No official API docs available; isolated stub only |
| E-commerce / online ordering | Not in scope |
| Accounting ledger / journal entries | Reports only, no double-entry bookkeeping |
| Email notifications | No email service in v1 |
| Mobile native app | Vue frontend is responsive; no native iOS/Android app |

---

## 6. Demo Story

> **Setup:** The database is seeded with a few categories (Beverages, Dairy,
> Snacks), ~20 products with barcodes, a supplier, and one user per role.

### Act 1 — Admin Sets the Stage

1. Admin logs in.
2. Admin creates a new cashier account "Marie".
3. Admin shows the user list with all five roles.

### Act 2 — Manager Prepares the Catalog

4. Manager logs in.
5. Manager adds a new product "Inyange Milk 500ml" with barcode `6789012345678`, cost 500 RWF, price 700 RWF, tax 18%.
6. Manager shows the product list, filters by "Dairy".

### Act 3 — Storekeeper Receives Stock

7. Storekeeper logs in.
8. Storekeeper creates a purchase order from supplier "Inyange Industries": 50× Inyange Milk.
9. Storekeeper marks the PO as received. Stock for Inyange Milk jumps to 50.

### Act 4 — The Wow Moment (Cashier Sale)

10. Cashier "Marie" logs in. The POS screen opens full-screen.
11. Marie scans barcode `6789012345678`. The product appears in the cart: "Inyange Milk 500ml — 700 RWF".
12. Marie scans two more products.
13. The screen shows subtotal, 18% VAT, and grand total.
14. Marie selects "Cash", enters amount tendered. The system calculates change.
15. Marie clicks **Complete Sale**.
16. A receipt appears with store name, TIN, items, tax breakdown, receipt number, and a mock EBM QR code.
17. Inventory for all three products has decreased by the sold quantities.

### Act 5 — Manager Reviews

18. Manager opens the dashboard. Today's revenue is updated. The sale appears in "Recent Sales".
19. Manager opens the Sales Report, filters by today, sees the sale details.
20. Manager checks inventory — Inyange Milk is now at 49.

### Epilogue — Accountant

21. Accountant logs in, views the sales report, exports it as CSV.
22. Accountant reprints the receipt from Act 4.

> **Total demo time target: 8–10 minutes.**

