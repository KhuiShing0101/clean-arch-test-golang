package extendloan

import (
	"errors"
	"library-management/internal/domain/loan"
)

// Custom errors
var (
	ErrLoanNotFound = errors.New("loan not found")
)

// ExtendLoanUseCase - Application layer use case for extending a loan's due date
type ExtendLoanUseCase struct {
	loanRepo loan.ILoanRepository
}

// NewExtendLoanUseCase creates a new use case instance
func NewExtendLoanUseCase(loanRepo loan.ILoanRepository) *ExtendLoanUseCase {
	return &ExtendLoanUseCase{
		loanRepo: loanRepo,
	}
}

// Execute processes the extend loan request
func (uc *ExtendLoanUseCase) Execute(req *ExtendLoanRequest) (*ExtendLoanResponse, error) {
	// 1. Find the loan
	loanEntity, err := uc.loanRepo.FindById(req.GetLoanId())
	if err != nil || loanEntity == nil {
		return nil, ErrLoanNotFound
	}

	// 2. Store original due date before extension
	originalDueDate := loanEntity.GetDueDate()

	// 3. Attempt to extend the loan
	// The ExtendDueDate method validates:
	// - Loan must be active (not returned)
	// - Loan must not be overdue
	// - Cannot exceed max extensions (2)
	if err := loanEntity.ExtendDueDate(); err != nil {
		return nil, err
	}

	// 4. Save the updated loan
	if err := uc.loanRepo.Save(loanEntity); err != nil {
		return nil, err
	}

	// 5. Create response with updated loan information
	response := NewExtendLoanResponse(
		string(loanEntity.GetId()),
		string(loanEntity.GetUserId()),
		string(loanEntity.GetBookId()),
		loanEntity.GetBorrowedAt(),
		originalDueDate,
		loanEntity.GetDueDate(),
		loanEntity.GetExtensionCount(),
	)

	return response, nil
}
