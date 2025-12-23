package book_test

import (
	"testing"
	"library-management/internal/domain/book"
)

func TestCanCreateValidISBN(t *testing.T) {
	isbn, err := book.NewISBN("9780134494166")
	if err != nil {
		t.Fatalf("Failed to create ISBN: %v", err)
	}

	if isbn.GetValue() != "9780134494166" {
		t.Errorf("Expected '9780134494166', got '%s'", isbn.GetValue())
	}
}

func TestCanCreateISBNWithHyphens(t *testing.T) {
	isbn, err := book.NewISBN("978-0-13-449416-6")
	if err != nil {
		t.Fatalf("Failed to create ISBN: %v", err)
	}

	if isbn.GetValue() != "9780134494166" {
		t.Errorf("Expected normalized '9780134494166', got '%s'", isbn.GetValue())
	}
}

func TestInvalidISBNReturnsError(t *testing.T) {
	_, err := book.NewISBN("invalid")
	if err == nil {
		t.Error("Expected error for invalid ISBN")
	}
}

func TestFormatsISBNCorrectly(t *testing.T) {
	isbn, _ := book.NewISBN("9780134494166")

	if isbn.GetFormatted() != "978-0-13-449416-6" {
		t.Errorf("Expected '978-0-13-449416-6', got '%s'", isbn.GetFormatted())
	}
}