CREATE DATABASE IF NOT EXISTS GuangShunCoinAction;
USE GuangShunCoinAction;

-- CREATE DATABASE `djangoDB`;
-- CREATE USER 'user' IDENTIFIED BY 'password';
-- GRANT ALL privileges ON `djangoDB`.* TO 'user'@'localhost';
-- FLUSH PRIVILEGES;

-- CREATE USER 'user'@'*' IDENTIFIED BY 'password';
-- GRANT ALL PRIVILEGES ON djangoDB.* TO 'user'@'*';
-- FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS User (
    user_id VARCHAR(36) PRIMARY KEY,
	username VARCHAR(36) NOT NULL,
    password VARCHAR(60) NOT NULL,
    cellphone VARCHAR(20),
    fb_account VARCHAR(100),
    email VARCHAR(100),
    address VARCHAR(255)
);
