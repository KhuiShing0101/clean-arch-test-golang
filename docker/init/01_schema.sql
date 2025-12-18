-- docker/init/01_schema.sql
-- Library Management System Schema

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status ENUM('ACTIVE', 'SUSPENDED', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    max_loans INT NOT NULL DEFAULT 5,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    -- ❌ NO status field! Single Source of Truth = loans table
);

CREATE TABLE IF NOT EXISTS loans (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    book_id VARCHAR(36) NOT NULL,
    borrowed_at DATETIME NOT NULL,
    due_date DATETIME NOT NULL,
    returned_at DATETIME DEFAULT NULL,  -- NULL = active loan (borrowed)
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id),
    INDEX idx_user_active (user_id, returned_at),
    INDEX idx_book_active (book_id, returned_at)
);