package returnbook

import (
	"time"

	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
)

// ReturnBookResponse - DTO for return book result
type ReturnBookResponse struct {
	LoanId      loan.LoanId
	BookId      book.BookId
	UserId      user.UserId
	BorrowedAt  time.Time
	DueDate     time.Time
	ReturnedAt  time.Time
	DaysLate    int
	LateFee     float64
	IsOverdue   bool
}

// NewReturnBookResponse creates a new response
func NewReturnBookResponse(
	loanId loan.LoanId,
	bookId book.BookId,
	userId user.UserId,
	borrowedAt time.Time,
	dueDate time.Time,
	returnedAt time.Time,
	daysLate int,
	lateFee float64,
	isOverdue bool,
) *ReturnBookResponse {
	return &ReturnBookResponse{
		LoanId:     loanId,
		BookId:     bookId,
		UserId:     userId,
		BorrowedAt: borrowedAt,
		DueDate:    dueDate,
		ReturnedAt: returnedAt,
		DaysLate:   daysLate,
		LateFee:    lateFee,
		IsOverdue:  isOverdue,
	}
}
