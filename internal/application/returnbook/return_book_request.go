package returnbook

import (
	"library-management/internal/domain/loan"
)

// ReturnBookRequest - DTO for returning a book
type ReturnBookRequest struct {
	LoanId loan.LoanId
}

// NewReturnBookRequest creates a new request
func NewReturnBookRequest(loanId loan.LoanId) *ReturnBookRequest {
	return &ReturnBookRequest{
		LoanId: loanId,
	}
}
