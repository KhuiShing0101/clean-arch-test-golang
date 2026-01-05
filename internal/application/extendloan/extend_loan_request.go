package extendloan

import (
	"errors"
	"library-management/internal/domain/loan"
)

// ExtendLoanRequest - DTO for extending a loan's due date
type ExtendLoanRequest struct {
	LoanId loan.LoanId
}

// NewExtendLoanRequest - Factory function with validation
func NewExtendLoanRequest(loanId string) (*ExtendLoanRequest, error) {
	if loanId == "" {
		return nil, errors.New("loan ID is required")
	}

	return &ExtendLoanRequest{
		LoanId: loan.LoanId(loanId),
	}, nil
}

// GetLoanId - Getter for loan ID
func (r *ExtendLoanRequest) GetLoanId() loan.LoanId {
	return r.LoanId
}
