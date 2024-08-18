CREATE DATABASE IF NOT EXISTS GuangShunCoinAction;
USE GuangShunCoinAction;

CREATE TABLE IF NOT EXISTS User (
    userId VARCHAR(36) PRIMARY KEY NOT NULL,
    userPasswd VARCHAR(60) NOT NULL,
    realName VARCHAR(36) NOT NULL,
    cellphone VARCHAR(20) NOT NULL,
    nickName VARCHAR(36) NOT NULL,
    -- username VARCHAR(36) NOT NULL,
    -- fbAccount VARCHAR(100),
    -- email VARCHAR(100),
    -- postcode VARCHAR(36),
    -- shippingAddr VARCHAR(255) NOT NULL,
    userRole VARCHAR(36) NOT NULL,
    loginStatus VARCHAR(36) NOT NULL,
    captcha VARCHAR(36)
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
    notified BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (userId) REFERENCES `User` (userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ProductImage (
    imageId VARCHAR(36) PRIMARY KEY NOT NULL,
    productId VARCHAR(36) NOT NULL,
    imageUrl VARCHAR(255) NOT NULL,
    seq BIGINT NOT NULL,
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

-- INSERT INTO User (userId, userPasswd, realName, cellphone, nickName, userRole, loginStatus) 
-- VALUES ('4dbe7ad6-776c-425a-8b78-73c9b031a9ae', '$2a$14$4IsnUa723JQbHrWy0GtZoOxenTg.HhI9TrsQvluDiCdiWXZcwNuiW', 'admin', '0912345678', 'admin', 'admin', 'signIn');
