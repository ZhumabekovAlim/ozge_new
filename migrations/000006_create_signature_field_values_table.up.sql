CREATE TABLE signature_field_values
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    signature_id      INT NOT NULL,
    contract_field_id INT NOT NULL,
    field_value       TEXT,
    FOREIGN KEY (signature_id) REFERENCES signatures (id) ON DELETE CASCADE,
    FOREIGN KEY (contract_field_id) REFERENCES contract_fields (id) ON DELETE CASCADE
);
