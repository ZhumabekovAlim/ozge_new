CREATE TABLE signatures
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    contract_id  INT NOT NULL,
    client_name  VARCHAR(255),
    client_iin   VARCHAR(12),
    client_phone VARCHAR(20),
    signed_at    TIMESTAMP,
    method      ENUM('sms', 'ecp') NOT NULL,
    FOREIGN KEY (contract_id) REFERENCES contracts (id) ON DELETE CASCADE
);
