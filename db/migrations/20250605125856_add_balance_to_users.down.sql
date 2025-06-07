-- Recreate accounts table (reverse of the up migration)
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    balance FLOAT NOT NULL DEFAULT 0,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Remove balance from users
ALTER TABLE users DROP COLUMN balance;
