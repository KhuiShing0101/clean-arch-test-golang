package controllers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"library-management/internal/application/getbook"
	"library-management/internal/http/controllers"
	queryImpl "library-management/internal/infrastructure/query"
)

func TestGetBookAPI_Returns200ForExistingBook(t *testing.T) {
	// Arrange - Set up test database
	db := setupTestDatabase(t)
	defer db.Close()

	// Seed database with test book
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author)
		VALUES ('b-12345', '9780134685991', 'Clean Architecture', 'Robert C. Martin')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Set up dependencies and controller
	deps := setupDependencies(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/books/{bookId}", deps.BookController.GetBook).Methods("GET")

	// Act - Make HTTP request
	req := httptest.NewRequest("GET", "/api/books/b-12345", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["success"] != true {
		t.Error("Expected success to be true")
	}

	data := response["data"].(map[string]interface{})
	if data["id"] != "b-12345" {
		t.Errorf("Expected book ID 'b-12345', got %v", data["id"])
	}
	if data["title"] != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got %v", data["title"])
	}
}

func TestGetBookAPI_Returns404ForNonExistentBook(t *testing.T) {
	// Arrange
	db := setupTestDatabase(t)
	defer db.Close()

	deps := setupDependencies(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/books/{bookId}", deps.BookController.GetBook).Methods("GET")

	// Act
	req := httptest.NewRequest("GET", "/api/books/b-99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["success"] != false {
		t.Error("Expected success to be false")
	}
	if response["error"] != "Book not found" {
		t.Errorf("Expected error 'Book not found', got %v", response["error"])
	}
}

// Dependencies holds test dependencies
type Dependencies struct {
	BookController *controllers.BookController
}

// setupDependencies wires up dependencies for testing
func setupDependencies(db *sql.DB) *Dependencies {
	bookQueryService := queryImpl.NewBookQueryService(db)
	getBookUseCase := getbook.NewGetBookUseCase(bookQueryService)
	bookController := controllers.NewBookController(getBookUseCase)

	return &Dependencies{
		BookController: bookController,
	}
}

func setupTestDatabase(t *testing.T) *sql.DB {
	// Connect to MySQL test database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/library")
	if err != nil {
		t.Fatal(err)
	}

	// Clean tables before each test
	cleanupTables(t, db)

	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		t.Fatal(err)
	}

	tables := []string{"loans", "books", "users"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		t.Fatal(err)
	}
}