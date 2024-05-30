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
    userId VARCHAR(36) PRIMARY KEY,
	username VARCHAR(36) NOT NULL,
    passwd VARCHAR(60) NOT NULL,
    realName VARCHAR(36),
    cellphone VARCHAR(20),
    fbAccount VARCHAR(100),
    email VARCHAR(100),
    postcode VARCHAR(36),
    shippingAddr VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS Product (
	productId INT PRIMARY KEY AUTO_INCREMENT,
    productName VARCHAR(36) NOT NULL,
    category VARCHAR(36),
    price DECIMAL,
    minBidPrice DECIMAL,
    imgUrl VARCHAR(255), 
    createTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    startTime TIMESTAMP,
    endedTime TIMESTAMP
    shippingStatus VARCHAR(36),
);


CREATE TABLE IF NOT EXISTS TrackingList (
    userId VARCHAR(36),
    productId VARCHAR(36),
    trackedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userId, productId),
    FOREIGN KEY (userId) REFERENCES User(userId) ON DELETE CASCADE,
    FOREIGN KEY (productId) REFERENCES Product(productId) ON DELETE CASCADE
);