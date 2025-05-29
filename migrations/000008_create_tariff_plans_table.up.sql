CREATE TABLE tariff_plans
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255)   NOT NULL,
    discount   DECIMAL(10, 2) NOT NULL,
    from_count INT            NOT NULL,
    to_count   INT            NOT NULL,
    is_active  BOOLEAN   DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);