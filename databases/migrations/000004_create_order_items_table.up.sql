CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    order_id UUID REFERENCES orders(id),
    product_id UUID REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);