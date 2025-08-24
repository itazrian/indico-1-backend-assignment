DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS settlements;
DROP TABLE IF EXISTS jobs;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    stock INT NOT NULL
);

CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT,
    buyer_id VARCHAR(50),
    quantity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    merchant_id VARCHAR(50),
    amount_cents INT,
    fee_cents INT,
    status VARCHAR(20),
    paid_at DATE
);

CREATE TABLE settlements (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    merchant_id VARCHAR(50),
    date DATE,
    gross_cents BIGINT,
    fee_cents BIGINT,
    net_cents BIGINT,
    txn_count INT,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unique_run_id VARCHAR(100),
    UNIQUE KEY uniq_merchant_date (merchant_id, date)
);

CREATE TABLE jobs (
    id VARCHAR(100) PRIMARY KEY,
    status VARCHAR(50),
    result_path VARCHAR(255),
    progress INT,
    processed BIGINT,
    total BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions (merchant_id, amount_cents, fee_cents, status, paid_at) VALUES
('mid1',1000,100,'PAID','2025-08-05'),
('mid1',2000,200,'PAID','2025-08-05'),
('mid2',1500,150,'PAID','2025-08-06'),
('mid1',1200,120,'PAID','2025-08-07'),
('mid2',1800,180,'PAID','2025-08-07');

INSERT INTO products (name, stock) VALUES ('T-Shirt', 300);

