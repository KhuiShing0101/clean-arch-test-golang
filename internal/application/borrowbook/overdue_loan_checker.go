package borrowbook

import (
	"fmt"
	"myapp/internal/domain/entity"
)

// OverdueLoanChecker is a domain service for checking overdue loans
//
// Single Responsibility: Verify no overdue loans exist
type OverdueLoanChecker struct{}

// NewOverdueLoanChecker creates a new OverdueLoanChecker
func NewOverdueLoanChecker() *OverdueLoanChecker {
	return &OverdueLoanChecker{}
}

// Check verifies if user has any overdue loans
//
// Returns error if overdue loans exist
func (c *OverdueLoanChecker) Check(overdueLoans []*entity.Loan) error {
	if len(overdueLoans) > 0 {
		return fmt.Errorf(
			"user has %d overdue loan(s): please return them before borrowing",
			len(overdueLoans),
		)
	}

	return nil
}