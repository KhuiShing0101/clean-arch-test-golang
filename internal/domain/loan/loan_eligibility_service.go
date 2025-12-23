package loan

import (
  "fmt"
  "library-management/internal/domain/user"
  "library-management/internal/domain/book"
)

// LoanEligibilityService - Domain Service for loan eligibility checks
type LoanEligibilityService struct{}

func NewLoanEligibilityService() *LoanEligibilityService {
  return &LoanEligibilityService{}
}

func (s *LoanEligibilityService) CanBorrow(u *user.User, b *book.Book) bool {
  // Rule 1: User must not have reached max loan limit
  if !u.CanBorrowMore() {
    return false
  }

  // Rule 2: User must not have overdue books
  if u.HasOverdueBooks() {
    return false
  }

  // Rule 3: Book must have available copies
  if !b.IsAvailable() {
    return false
  }

  return true
}

func (s *LoanEligibilityService) GetIneligibilityReason(u *user.User, b *book.Book) *string {
  if !u.CanBorrowMore() {
    reason := fmt.Sprintf("User has reached maximum loan limit (%d books)", u.GetMaxLoans())
    return &reason
  }

  if u.HasOverdueBooks() {
    reason := "User has overdue books"
    return &reason
  }

  if !b.IsAvailable() {
    reason := fmt.Sprintf("No copies available for book \"%s\"", b.GetTitle())
    return &reason
  }

  return nil
}