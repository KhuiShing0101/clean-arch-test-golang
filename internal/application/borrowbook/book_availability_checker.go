package borrowbook

import (
	"context"
	"fmt"
	"myapp/internal/domain/entity"
)

// BookAvailabilityChecker is a domain service for checking book availability
//
// Single Source of Truth: Book availability is determined by
// checking if an active loan exists (not by book.status!)
type BookAvailabilityChecker struct {
	loanRepository LoanRepositoryInterface
}

// NewBookAvailabilityChecker creates a new BookAvailabilityChecker
func NewBookAvailabilityChecker(loanRepo LoanRepositoryInterface) *BookAvailabilityChecker {
	return &BookAvailabilityChecker{
		loanRepository: loanRepo,
	}
}

// Check verifies if book is available for borrowing
//
// Returns error if book is not available
func (c *BookAvailabilityChecker) Check(ctx context.Context, book *entity.Book) error {
	// Single Source of Truth: Check loans table, NOT book.status!
	activeLoan, err := c.loanRepository.FindActiveLoanByBookID(ctx, book.GetID())
	if err != nil {
		return fmt.Errorf("failed to check book availability: %w", err)
	}

	if activeLoan != nil {
		return fmt.Errorf("book is not available: already borrowed")
	}

	return nil
}