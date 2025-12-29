package borrowbook

import (
	"fmt"
	"myapp/internal/domain/entity"
)

// LoanLimitChecker is a domain service for checking loan count limits
//
// Single Responsibility: Verify user hasn't exceeded borrowing limit
type LoanLimitChecker struct {
	loanGateway LoanGatewayInterface
}

// NewLoanLimitChecker creates a new LoanLimitChecker
func NewLoanLimitChecker(loanGateway LoanGatewayInterface) *LoanLimitChecker {
	return &LoanLimitChecker{
		loanGateway: loanGateway,
	}
}

// Check verifies if user is within borrowing limit
//
// Returns error if limit exceeded
func (c *LoanLimitChecker) Check(user *entity.User, activeLoans []*entity.Loan) error {
	maxLoans := user.GetMaxLoans() // From User entity (5 books)
	currentCount := len(activeLoans)

	if currentCount >= maxLoans {
		return fmt.Errorf(
			"loan limit exceeded: user has %d active loans (max: %d)",
			currentCount,
			maxLoans,
		)
	}

	return nil
}