CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    order_id UUID REFERENCES orders(id),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    transaction_id int default serial,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
