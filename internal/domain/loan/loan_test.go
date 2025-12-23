package loan_test

import (
	"testing"
	"time"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
	"library-management/internal/domain/book"
)

func TestDueDateIsCalculatedCorrectly(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanId, _ := loan.NewLoanId("loan-123")
	userId, _ := user.NewUserId("user-456")
	bookId, _ := book.NewBookId("book-789")

	l := loan.NewLoan(loanId, userId, bookId, borrowedAt, nil)

	expectedDueDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	if !l.GetDueDate().Equal(expectedDueDate) {
		t.Errorf("Expected due date %v, got %v", expectedDueDate, l.GetDueDate())
	}
}

func TestLoanIsOverdueAfterDueDate(t *testing.T) {
	borrowedAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create Value Objects
	loanId, _ := loan.NewLoanId("loan-123")
	userId, _ := user.NewUserId("user-456")
	bookId, _ := book.NewBookId("book-789")

	l := loan.NewLoan(loanId, userId, bookId, borrowedAt, nil)

	if !l.IsOverdue(&currentDate) {
		t.Error("Expected loan to be overdue")
	}
}

func TestCanMarkLoanAsReturned(t *testing.T) {
	// Create Value Objects
	loanId, _ := loan.NewLoanId("loan-123")
	userId, _ := user.NewUserId("user-456")
	bookId, _ := book.NewBookId("book-789")

	l := loan.NewLoan(loanId, userId, bookId, time.Now(), nil)

	if l.IsReturned() {
		t.Error("Loan should not be returned initially")
	}

	returnedAt := time.Now()
	err := l.MarkAsReturned(&returnedAt)
	if err != nil {
		t.Fatalf("Failed to mark as returned: %v", err)
	}

	if !l.IsReturned() {
		t.Error("Loan should be marked as returned")
	}
	if l.GetReturnedAt() == nil || l.GetReturnedAt().IsZero() {
		t.Error("ReturnedAt should be set")
	}
}
