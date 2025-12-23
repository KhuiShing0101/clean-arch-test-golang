package loan_test

import (
	"testing"
	"time"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
	"library-management/internal/domain/book"
)

func TestLoanEligibility(t *testing.T) {
	t.Run("UserCanBorrow_AllConditionsMet", func(t *testing.T) {
		service := loan.NewLoanEligibilityService()

		// Create user (Part 2: NewUser with name and email)
		u := user.NewUser("John Doe", "john@example.com")

		// Create book
		bookId, _ := book.NewBookId("book-456")
		isbn, _ := book.NewISBN("9780134494166")
		b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

		if !service.CanBorrow(u, b) {
			t.Error("User should be able to borrow")
		}
	})

	t.Run("UserCannotBorrow_MaxLoansReached", func(t *testing.T) {
		service := loan.NewLoanEligibilityService()

		// Reconstruct user with max borrow count
		userId, _ := user.NewUserId("12345678")
		u := user.ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			user.UserStatusActive,
			user.MaxBorrowLimit, // At limit
			0,
			time.Now(),
		)

		// Create book
		bookId, _ := book.NewBookId("book-456")
		isbn, _ := book.NewISBN("9780134494166")
		b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

		if service.CanBorrow(u, b) {
			t.Error("User should not be able to borrow (max loans reached)")
		}

		// Check reason
		reason := service.GetIneligibilityReason(u, b)
		if reason == nil {
			t.Error("Expected ineligibility reason")
		}
	})

	t.Run("UserCannotBorrow_Suspended", func(t *testing.T) {
		service := loan.NewLoanEligibilityService()

		// Reconstruct suspended user
		userId, _ := user.NewUserId("12345678")
		u := user.ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			user.UserStatusSuspended,
			0,
			0,
			time.Now(),
		)

		// Create book
		bookId, _ := book.NewBookId("book-456")
		isbn, _ := book.NewISBN("9780134494166")
		b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

		if service.CanBorrow(u, b) {
			t.Error("Suspended user should not be able to borrow")
		}

		// Check reason
		reason := service.GetIneligibilityReason(u, b)
		if reason == nil {
			t.Error("Expected ineligibility reason for suspended user")
		}
	})

	t.Run("UserCannotBorrow_OverdueFees", func(t *testing.T) {
		service := loan.NewLoanEligibilityService()

		// Reconstruct user with overdue fees
		userId, _ := user.NewUserId("12345678")
		u := user.ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			user.UserStatusActive,
			0,
			10.50, // Has overdue fees
			time.Now(),
		)

		// Create book
		bookId, _ := book.NewBookId("book-456")
		isbn, _ := book.NewISBN("9780134494166")
		b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 3)

		if service.CanBorrow(u, b) {
			t.Error("User with overdue fees should not be able to borrow")
		}

		// Check reason
		reason := service.GetIneligibilityReason(u, b)
		if reason == nil {
			t.Error("Expected ineligibility reason for overdue fees")
		}
	})

	t.Run("UserCannotBorrow_NoAvailableCopies", func(t *testing.T) {
		service := loan.NewLoanEligibilityService()

		// Create active user
		u := user.NewUser("John Doe", "john@example.com")

		// Create book with 1 copy, then borrow it to make unavailable
		bookId, _ := book.NewBookId("book-456")
		isbn, _ := book.NewISBN("9780134494166")
		b, _ := book.NewBook(bookId, "Clean Architecture", "Robert Martin", isbn, 1)
		b.BorrowCopy() // Now availableCopies = 0

		if service.CanBorrow(u, b) {
			t.Error("User should not be able to borrow (no available copies)")
		}

		// Check reason
		reason := service.GetIneligibilityReason(u, b)
		if reason == nil {
			t.Error("Expected ineligibility reason for no copies")
		}
	})
}
