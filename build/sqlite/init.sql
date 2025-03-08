-- init.sql for SQLite
-- Note: SQLite doesn't use CREATE DATABASE or USE statements

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert user data
INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Wilson', 'bob@example.com');

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    price REAL,
    stock INTEGER
);

-- Insert product data
INSERT INTO products (name, price, stock) VALUES
    ('Laptop', 999.99, 10),
    ('Mouse', 24.99, 50),
    ('Keyboard', 59.99, 30);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_amount REAL NOT NULL CHECK (total_amount >= 0),
    user_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    
    -- Unique constraint
    UNIQUE (id, user_id)
);

-- Create trigger to update the updated_at timestamp when a row is updated
CREATE TRIGGER IF NOT EXISTS update_orders_timestamp 
AFTER UPDATE ON orders
FOR EACH ROW
BEGIN
    UPDATE orders SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Insert order data
INSERT INTO orders (user_id, total_amount, status) VALUES
    (1, 299.99, 'paid'),
    (1, 599.99, 'pending'),
    (2, 1299.99, 'paid'),
    (2, 49.99, 'cancelled'),
    (3, 799.99, 'paid'),
    (1, 99.99, 'refunded'),
    (3, 199.99, 'pending'),
    (2, 899.99, 'paid'),
    (1, 149.99, 'pending'),
    (3, 399.99, 'cancelled');
