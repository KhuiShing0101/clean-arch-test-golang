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
  // Part 2: User.CanBorrow() checks status, limit, AND fees
  if !u.CanBorrow() {
    return false
  }

  // Rule: Book must have available copies
  if !b.IsAvailable() {
    return false
  }

  return true
}

func (s *LoanEligibilityService) GetIneligibilityReason(u *user.User, b *book.Book) *string {
  if u.Status() == user.UserStatusSuspended {
    reason := "User is suspended"
    return &reason
  }

  if u.CurrentBorrowCount() >= user.MaxBorrowLimit {
    reason := fmt.Sprintf("User has reached maximum loan limit (%d books)", user.MaxBorrowLimit)
    return &reason
  }

  if u.OverdueFees() > 0 {
    reason := fmt.Sprintf("User has overdue fees: $%.2f", u.OverdueFees())
    return &reason
  }

  if !b.IsAvailable() {
    reason := fmt.Sprintf("No copies available for book \"%s\"", b.GetTitle())
    return &reason
  }

  return nil
}