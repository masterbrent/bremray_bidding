-- Create items table
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price >= 0),
    category VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on category for filtering
CREATE INDEX idx_items_category ON items(category);

-- Create index on created_at for sorting
CREATE INDEX idx_items_created_at ON items(created_at DESC);