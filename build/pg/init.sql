-- init.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Wilson', 'bob@example.com');

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    price DECIMAL(10,2),
    stock INTEGER
);

INSERT INTO products (name, price, stock) VALUES
    ('Laptop', 999.99, 10),
    ('Mouse', 24.99, 50),
    ('Keyboard', 59.99, 30);

    CREATE TABLE orders (
   id SERIAL PRIMARY KEY,
   user_id INTEGER NOT NULL,
   total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
   status VARCHAR(20) NOT NULL DEFAULT 'pending',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   
   -- Foreign key constraint
   CONSTRAINT fk_user_id 
       FOREIGN KEY (user_id) 
       REFERENCES users(id)
       ON DELETE RESTRICT
       ON UPDATE CASCADE,
   
   -- Check constraint
   CONSTRAINT valid_status 
       CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
   
   -- Unique constraint
   CONSTRAINT unique_order_reference 
       UNIQUE (id, user_id)
);

-- Insert orders
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
