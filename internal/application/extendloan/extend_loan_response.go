package extendloan

import (
	"time"
)

// ExtendLoanResponse - DTO for extended loan information
type ExtendLoanResponse struct {
	LoanId         string
	UserId         string
	BookId         string
	BorrowedAt     time.Time
	OriginalDueDate time.Time
	NewDueDate     time.Time
	ExtensionCount int
	Message        string
}

// NewExtendLoanResponse - Factory function for creating response
func NewExtendLoanResponse(
	loanId string,
	userId string,
	bookId string,
	borrowedAt time.Time,
	originalDueDate time.Time,
	newDueDate time.Time,
	extensionCount int,
) *ExtendLoanResponse {
	return &ExtendLoanResponse{
		LoanId:         loanId,
		UserId:         userId,
		BookId:         bookId,
		BorrowedAt:     borrowedAt,
		OriginalDueDate: originalDueDate,
		NewDueDate:     newDueDate,
		ExtensionCount: extensionCount,
		Message:        "Loan extended successfully",
	}
}
