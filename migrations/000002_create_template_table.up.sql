CREATE TABLE templates
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT          NOT NULL,
    name       VARCHAR(255) NOT NULL,
    file_path  VARCHAR(500) NOT NULL, -- только обработанный шаблон
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);
