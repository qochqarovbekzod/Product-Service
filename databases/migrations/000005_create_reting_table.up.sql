CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    product_id UUID REFERENCES products(id),
    user_id UUID NOT NULL,
    rating DECIMAL(2, 1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
