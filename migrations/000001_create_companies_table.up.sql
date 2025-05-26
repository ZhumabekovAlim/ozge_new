CREATE TABLE companies
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    phone        VARCHAR(100) UNIQUE,
    password    VARCHAR(100),
    email        VARCHAR(100) UNIQUE
);

