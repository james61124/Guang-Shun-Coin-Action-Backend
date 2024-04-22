CREATE DATABASE IF NOT EXISTS djangoDB;
USE djangoDB;

-- CREATE DATABASE `djangoDB`;
-- CREATE USER 'user' IDENTIFIED BY 'password';
-- GRANT ALL privileges ON `djangoDB`.* TO 'user'@'localhost';
-- FLUSH PRIVILEGES;

-- CREATE USER 'user'@'*' IDENTIFIED BY 'password';
-- GRANT ALL PRIVILEGES ON djangoDB.* TO 'user'@'*';
-- FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS User (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    account VARCHAR(36) UNIQUE,
	username VARCHAR(36),
    password VARCHAR(60) NOT NULL,
    cellphone VARCHAR(20),
    fb_account VARCHAR(100),
    email VARCHAR(100),
    address VARCHAR(255)
);
