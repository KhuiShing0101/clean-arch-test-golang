package loan

import (
  "errors"
  "time"
  "library-management/internal/domain/user"
  "library-management/internal/domain/book"
)

// Loan Entity - represents a book loan
type Loan struct {
  id         *LoanId
  userId     *user.UserId
  bookId     *book.BookId
  borrowedAt time.Time
  dueDate    time.Time
  returnedAt *time.Time
}

const LoanPeriodDays = 14

func NewLoan(
  id *LoanId,
  userId *user.UserId,
  bookId *book.BookId,
  borrowedAt time.Time,
  returnedAt *time.Time,
) *Loan {
  if borrowedAt.IsZero() {
    borrowedAt = time.Now()
  }

  dueDate := calculateDueDate(borrowedAt)

  return &Loan{
    id:         id,
    userId:     userId,
    bookId:     bookId,
    borrowedAt: borrowedAt,
    dueDate:    dueDate,
    returnedAt: returnedAt,
  }
}

func calculateDueDate(borrowedAt time.Time) time.Time {
  return borrowedAt.AddDate(0, 0, LoanPeriodDays)
}

func (l *Loan) IsOverdue(currentDate *time.Time) bool {
  var now time.Time
  if currentDate != nil {
    now = *currentDate
  } else {
    now = time.Now()
  }
  return now.After(l.dueDate) && l.returnedAt == nil
}

func (l *Loan) IsReturned() bool {
  return l.returnedAt != nil
}

func (l *Loan) MarkAsReturned(returnedAt *time.Time) error {
  if l.IsReturned() {
    return errors.New("loan has already been returned")
  }

  if returnedAt != nil {
    l.returnedAt = returnedAt
  } else {
    now := time.Now()
    l.returnedAt = &now
  }

  return nil
}

func (l *Loan) GetDaysUntilDue(currentDate *time.Time) int {
  var now time.Time
  if currentDate != nil {
    now = *currentDate
  } else {
    now = time.Now()
  }

  duration := l.dueDate.Sub(now)
  days := int(duration.Hours() / 24)

  return days
}

// Getters
func (l *Loan) GetId() *LoanId {
  return l.id
}

func (l *Loan) GetUserId() *user.UserId {
  return l.userId
}

func (l *Loan) GetBookId() *book.BookId {
  return l.bookId
}

func (l *Loan) GetBorrowedAt() time.Time {
  return l.borrowedAt
}

func (l *Loan) GetDueDate() time.Time {
  return l.dueDate
}

func (l *Loan) GetReturnedAt() *time.Time {
  return l.returnedAt
}

func (l *Loan) GetLoanPeriodDays() int {
  return LoanPeriodDays
}