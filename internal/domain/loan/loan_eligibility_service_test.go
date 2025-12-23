package loan_test

import (
	"testing"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
	"library-management/internal/domain/book"
)

func TestUserCanBorrowWhenAllConditionsMet(t *testing.T) {
	service := loan.NewLoanEligibilityService()

	// Create Value Objects
	userId, _ := user.NewUserId("user-123")
	u, _ := user.NewUser(userId, "John Doe", "john@example.com", 0, false)

	bookId, _ := book.NewBookId("book-456")
	isbn, _ := book.NewISBN("9780134494166")
	b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

	if !service.CanBorrow(u, b) {
		t.Error("User should be able to borrow")
	}
}

func TestUserCannotBorrowWhenMaxLoansReached(t *testing.T) {
	service := loan.NewLoanEligibilityService()

	// Create Value Objects
	userId, _ := user.NewUserId("user-123")
	u, _ := user.NewUser(userId, "John Doe", "john@example.com", 5, false)

	bookId, _ := book.NewBookId("book-456")
	isbn, _ := book.NewISBN("9780134494166")
	b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

	if service.CanBorrow(u, b) {
		t.Error("User should not be able to borrow (max loans reached)")
	}
}
