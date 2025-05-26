CREATE TABLE contract_fields
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    contract_id INT          NOT NULL,
    field_name  VARCHAR(255) NOT NULL,
    field_type  VARCHAR(50) DEFAULT 'text',
    FOREIGN KEY (contract_id) REFERENCES contracts (id) ON DELETE CASCADE
);
