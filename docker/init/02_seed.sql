-- docker/init/02_seed.sql
-- Sample Data for Testing

-- Users
INSERT INTO users (id, name, email, status, max_loans) VALUES
('user-001', '田中太郎', 'tanaka@example.com', 'ACTIVE', 5),
('user-002', '佐藤花子', 'sato@example.com', 'ACTIVE', 5),
('user-003', '鈴木一郎', 'suzuki@example.com', 'SUSPENDED', 5),
('user-004', '高橋美咲', 'takahashi@example.com', 'ACTIVE', 3);

-- Books
INSERT INTO books (id, title, author, isbn) VALUES
('book-001', 'Clean Code', 'Robert C. Martin', '9780132350884'),
('book-002', 'Clean Architecture', 'Robert C. Martin', '9780134494166'),
('book-003', 'Domain-Driven Design', 'Eric Evans', '9780321125217'),
('book-004', 'Refactoring', 'Martin Fowler', '9780134757599'),
('book-005', 'The Pragmatic Programmer', 'David Thomas', '9780135957059');

-- Loans (some active, some returned, some overdue)
INSERT INTO loans (id, user_id, book_id, borrowed_at, due_date, returned_at) VALUES
-- Active loan (book-001 is borrowed)
('loan-001', 'user-001', 'book-001', NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY), NULL),
-- Returned loan
('loan-002', 'user-001', 'book-002', DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 16 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY)),
-- Overdue loan (book-003 is overdue)
('loan-003', 'user-002', 'book-003', DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), NULL);