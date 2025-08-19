-- Create database
CREATE DATABASE IF NOT EXISTS ordersdb;
USE ordersdb;

-- Create locations table
CREATE TABLE IF NOT EXISTS locations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    UNIQUE(name, latitude, longitude)
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    orderId INT AUTO_INCREMENT PRIMARY KEY,
    resLocationId INT NOT NULL,
    cusLocationId INT NOT NULL,
    prepTimeInMinutes DOUBLE NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resLocationId) REFERENCES locations(id),
    FOREIGN KEY (cusLocationId) REFERENCES locations(id)
);

CREATE INDEX idx_orders_resLocationId ON orders (resLocationId);
CREATE INDEX idx_orders_cusLocationId ON orders (cusLocationId);