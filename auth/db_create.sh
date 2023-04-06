#!/bin/sh
CREATE TABLE your_table_name (
  id INT AUTO_INCREMENT PRIMARY KEY,
  openid VARCHAR(64) NOT NULL UNIQUE,
  accesstoken VARCHAR(128) NOT NULL,
  name VARCHAR(64) NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX (openid)
);
