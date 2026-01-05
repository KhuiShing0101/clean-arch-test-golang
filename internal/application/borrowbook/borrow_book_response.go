package borrowbook

import (
	"time"

	"library-management/internal/domain/loan"
)

type BorrowBookResponse struct {
	LoanId     string
	UserId     string
	BookId     string
	BookTitle  string
	BorrowedAt time.Time
	DueDate    time.Time
}

func NewBorrowBookResponse(
	l *loan.Loan,
	bookTitle string,
) *BorrowBookResponse {
	return &BorrowBookResponse{
		LoanId:     string(l.GetId()),
		UserId:     string(l.GetUserId()),
		BookId:     string(l.GetBookId()),
		BookTitle:  bookTitle,
		BorrowedAt: l.GetBorrowedAt(),
		DueDate:    l.GetDueDate(),
	}
}