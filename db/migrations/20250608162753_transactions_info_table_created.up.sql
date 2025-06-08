
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    sender_id INTEGER,
    receiver_id INTEGER,
    amount DECIMAL(10,2),  -- or NUMERIC(10,2) for precise money handling
    type VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);