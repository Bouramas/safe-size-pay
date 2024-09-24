CREATE DATABASE IF NOT EXISTS safe_size_db;
USE safe_size_db;

CREATE TABLE IF NOT EXISTS Users (
    id              BINARY(16)      NOT NULL,
    name            VARCHAR(100)    NOT NULL,
    email           VARCHAR(100)    NOT NULL UNIQUE,
    password_hash   VARCHAR(255)    NOT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;


CREATE TABLE IF NOT EXISTS Transactions(
     id           BINARY(16) NOT NULL,
     user_id      BINARY(16) NOT NULL,
     order_id     INT NULL,
     order_msg    VARCHAR(500) NULL,
     description  TEXT,
     amount       DECIMAL(10, 2) NOT NULL,
     order_status ENUM('pending', 'success', 'failed') NOT NULL DEFAULT 'pending',
     created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
     PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;