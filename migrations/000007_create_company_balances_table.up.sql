CREATE TABLE company_balances
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    company_id        INT NOT NULL UNIQUE,
    sms_signatures    INT DEFAULT 5,
    ecp_signatures    INT DEFAULT 2,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);