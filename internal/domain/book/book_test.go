package book_test

import (
	"testing"
	"library-management/internal/domain/book"
)

func TestCanCreateBook(t *testing.T) {
	// Create Value Objects
	bookId, _ := book.NewBookId("book-123")
	isbn, _ := book.NewISBN("9780134494166")

	b, err := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	if b.GetId().GetValue() != "book-123" {
		t.Errorf("Expected ID 'book-123', got '%s'", b.GetId().GetValue())
	}
	if b.GetTitle() != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got '%s'", b.GetTitle())
	}
	if b.GetAvailableCopies() != 3 {
		t.Errorf("Expected 3 copies, got %d", b.GetAvailableCopies())
	}
}

func TestCanBorrowWhenCopiesAvailable(t *testing.T) {
	bookId, _ := book.NewBookId("book-123")
	isbn, _ := book.NewISBN("9780134494166")
	b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 2)

	if !b.IsAvailable() {
		t.Error("Expected book to be available for borrowing")
	}
}

func TestBorrowDecreasesAvailableCopies(t *testing.T) {
	bookId, _ := book.NewBookId("book-123")
	isbn, _ := book.NewISBN("9780134494166")
	b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 2)

	err := b.BorrowCopy()
	if err != nil {
		t.Fatalf("Borrow failed: %v", err)
	}

	if b.GetAvailableCopies() != 1 {
		t.Errorf("Expected 1 copy remaining, got %d", b.GetAvailableCopies())
	}
}

func TestCannotBorrowWhenNoCopies(t *testing.T) {
	bookId, _ := book.NewBookId("book-123")
	isbn, _ := book.NewISBN("9780134494166")
	b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 1)

	// Borrow the only copy
	_ = b.BorrowCopy()

	// Try to borrow again when no copies available
	err := b.BorrowCopy()
	if err == nil {
		t.Error("Expected error when borrowing with 0 available copies")
	}
}
