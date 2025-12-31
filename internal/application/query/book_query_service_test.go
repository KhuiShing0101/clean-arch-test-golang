package query_test

import (
	"database/sql"
	"testing"

	"library-management/internal/infrastructure/query"
	_ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/library")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean tables before each test
	cleanupTables(t, db)

	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	// Disable foreign key checks temporarily
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		t.Fatalf("Failed to disable foreign key checks: %v", err)
	}

	// Truncate tables in reverse dependency order
	tables := []string{"loans", "books", "users"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}

	// Re-enable foreign key checks
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		t.Fatalf("Failed to enable foreign key checks: %v", err)
	}
}

func TestBookQueryService_GetBookById_FindsAvailableBook(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	// Insert test book
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '9780123456789', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	queryService := query.NewBookQueryService(db)

	// Act
	book, err := queryService.GetBookById("b-12345")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.ID != "b-12345" {
		t.Errorf("Expected ID 'b-12345', got: %s", book.ID)
	}
	if book.Title != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got: %s", book.Title)
	}
	if !book.IsAvailable {
		t.Error("Expected book to be available")
	}
	if book.CurrentLoan != nil {
		t.Error("Expected no current loan for available book")
	}
}

func TestBookQueryService_GetBookById_FindsBorrowedBook(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	// Insert test book and loan
	// Insert user first (required by foreign key)
	_, err := db.Exec(`
		INSERT INTO users (id, name, email, status, created_at)
		VALUES ('u-001', 'Test User', 'test@example.com', 'ACTIVE', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '9780123456789', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO loans (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ('l-001', 'b-12345', 'u-001', '2024-01-15 10:00:00',
				'2024-02-15 10:00:00', NULL)
	`)
	if err != nil {
		t.Fatalf("Failed to insert loan: %v", err)
	}

	queryService := query.NewBookQueryService(db)

	// Act
	book, err := queryService.GetBookById("b-12345")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.IsAvailable {
		t.Error("Expected book to be borrowed (not available)")
	}
	if book.CurrentLoan == nil {
		t.Fatal("Expected current loan to exist")
	}
	if book.CurrentLoan.LoanID != "l-001" {
		t.Errorf("Expected loan ID 'l-001', got: %s", book.CurrentLoan.LoanID)
	}
}

func TestBookQueryService_GetBookById_ReturnsNilWhenNotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	queryService := query.NewBookQueryService(db)

	// Act
	book, err := queryService.GetBookById("b-99999")

	// Assert
	if err != nil {
		t.Errorf("Expected no error for not found case, got: %v", err)
	}
	if book != nil {
		t.Error("Expected nil for non-existent book")
	}
}