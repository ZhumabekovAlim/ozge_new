CREATE TABLE contracts
(
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    company_id          INT                 NOT NULL,
    template_id         INT                 NOT NULL,
    contract_token      VARCHAR(100) UNIQUE NOT NULL,
    generated_file_path VARCHAR(500), -- готовый файл для клиента
    client_filled       BOOLEAN   DEFAULT FALSE,
    company_sign        BOOLEAN   DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies (id),
    FOREIGN KEY (template_id) REFERENCES templates (id)
);
