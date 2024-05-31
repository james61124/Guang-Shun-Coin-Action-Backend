CREATE DATABASE IF NOT EXISTS GuangShunCoinAction;
USE GuangShunCoinAction;

-- CREATE DATABASE `djangoDB`;
-- CREATE USER 'user' IDENTIFIED BY 'password';
-- GRANT ALL privileges ON `djangoDB`.* TO 'user'@'localhost';
-- FLUSH PRIVILEGES;

-- CREATE USER 'user'@'*' IDENTIFIED BY 'password';
-- GRANT ALL PRIVILEGES ON djangoDB.* TO 'user'@'*';
-- FLUSH PRIVILEGES;

-- CREATE TABLE IF NOT EXISTS User (
--     userId VARCHAR(36) PRIMARY KEY,
-- 	username VARCHAR(36) NOT NULL,
--     passwd VARCHAR(60) NOT NULL,
--     realName VARCHAR(36),
--     cellphone VARCHAR(20),
--     fbAccount VARCHAR(100),
--     email VARCHAR(100),
--     postcode VARCHAR(36),
--     shippingAddr VARCHAR(255)
-- );

-- CREATE TABLE IF NOT EXISTS Product (
-- 	productId VARCHAR PRIMARY KEY,
--     productName VARCHAR(36) NOT NULL,
--     category VARCHAR(36),
--     price DECIMAL,
--     minBidPrice DECIMAL,
--     imgUrl VARCHAR(255), 
--     createTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     startTime TIMESTAMP,
--     endedTime TIMESTAMP
--     shippingStatus VARCHAR(36),
-- );



CREATE TABLE IF NOT EXISTS User (
    userId VARCHAR(36) PRIMARY KEY NOT NULL,
    username VARCHAR(36) NOT NULL,
    userPasswd VARCHAR(60) NOT NULL,
    realName VARCHAR(36) NOT NULL,
    cellphone VARCHAR(20) NOT NULL,
    fbAccount VARCHAR(100),
    email VARCHAR(100),
    postcode VARCHAR(36),
    shippingAddr VARCHAR(255) NOT NULL,
    userRole VARCHAR(36) NOT NULL,
    loginStatus VARCHAR(36) NOT NULL,
    nickName VARCHAR(36) NOT NULL
);

CREATE TABLE IF NOT EXISTS Product (
    productId VARCHAR(36) PRIMARY KEY NOT NULL,
    userId VARCHAR(36) NOT NULL,   
    productName VARCHAR(36) NOT NULL,
    category VARCHAR(36) NOT NULL,
    price BIGINT NOT NULL,
    minBidPrice BIGINT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    endedAt TIMESTAMP DEFAULT NULL,
    startAt TIMESTAMP NOT NULL,
    shippingStatus VARCHAR(36) NOT NULL,
    productDescription VARCHAR(255) NOT NULL,
    FOREIGN KEY (userId) REFERENCES `User` (userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ProductImage (
    imageId VARCHAR(36) PRIMARY KEY NOT NULL,
    productId VARCHAR(36) NOT NULL,
    imageUrl VARCHAR(255) NOT NULL,
    FOREIGN KEY (productId) REFERENCES `Product` (productId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS TrackingList (
    trackingId VARCHAR(36) PRIMARY KEY NOT NULL,
    userId VARCHAR(36) NOT NULL,
    productId VARCHAR(36) NOT NULL,
    FOREIGN KEY (userId) REFERENCES `User`(userId) ON DELETE CASCADE,
    FOREIGN KEY (productId) REFERENCES `Product`(productId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS History (
    historyId VARCHAR(36) PRIMARY KEY NOT NULL,
    productId VARCHAR(36) NOT NULL,
    userId VARCHAR(36) NOT NULL,
    bidPrice BIGINT NOT NULL,
    bidTime TIMESTAMP NOT NULL,
    status VARCHAR(36),
    FOREIGN KEY (productId) REFERENCES `Product`(productId) ON DELETE CASCADE,
    FOREIGN KEY (userId) REFERENCES `User`(userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CustomService (
    messageId VARCHAR(36) PRIMARY KEY NOT NULL,
    userId VARCHAR(36) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content LONGTEXT NOT NULL,
    FOREIGN KEY (userId) REFERENCES `User`(userId) ON DELETE CASCADE
);
