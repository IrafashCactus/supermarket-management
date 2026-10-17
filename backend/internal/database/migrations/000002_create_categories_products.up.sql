CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    barcode VARCHAR(100) NOT NULL UNIQUE,
    sku VARCHAR(100),
    category_id UUID NOT NULL REFERENCES categories(id),
    unit_price INTEGER NOT NULL,
    cost_price INTEGER NOT NULL,
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 18.00,
    unit_of_measure VARCHAR(50) NOT NULL DEFAULT 'piece',
    reorder_level INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
