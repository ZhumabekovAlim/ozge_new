CREATE TABLE payment_requests
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    company_id        INT NOT NULL,
    tariff_plan_id    INT,
    sms_count         INT DEFAULT 0,
    ecp_count         INT DEFAULT 0,
    total_amount      DECIMAL(10, 2) NOT NULL,
    status            ENUM('pending', 'paid', 'cancelled', 'expired') DEFAULT 'pending',
    payment_url       VARCHAR(500),
    payment_ref       VARCHAR(255),
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at           TIMESTAMP NULL,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    FOREIGN KEY (tariff_plan_id) REFERENCES tariff_plans (id) ON DELETE SET NULL
);
