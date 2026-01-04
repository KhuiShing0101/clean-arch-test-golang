package loan

import (
	"errors"
	"time"

	"library-management/internal/domain/book"
	"library-management/internal/domain/user"
	"github.com/google/uuid"
)

// Value Objects
type LoanId string

// Helper function to generate unique loan ID
func GenerateLoanId() LoanId {
	return LoanId(uuid.New().String())
}

const LoanPeriodDays = 14

type Loan struct {
	id         LoanId
	userId     user.UserId
	bookId     book.BookId
	borrowedAt time.Time
	dueDate    time.Time
	returnedAt *time.Time
}

// NewLoan - creates new loan (used in BorrowBook)
func NewLoan(userId user.UserId, bookId book.BookId) *Loan {
	now := time.Now()
	dueDate := now.AddDate(0, 0, LoanPeriodDays)

	return &Loan{
		id:         GenerateLoanId(),
		userId:     userId,
		bookId:     bookId,
		borrowedAt: now,
		dueDate:    dueDate,
		returnedAt: nil, // Not returned yet
	}
}

// IsActive - checks if loan is active
func (l *Loan) IsActive() bool {
	return l.returnedAt == nil
}

// RecordReturn - records return (used in ReturnBook)
func (l *Loan) RecordReturn() error {
	if !l.IsActive() {
		return errors.New("loan is already returned")
	}
	now := time.Now()
	l.returnedAt = &now
	return nil
}

// Getters
func (l *Loan) GetId() LoanId {
	return l.id
}

func (l *Loan) GetUserId() user.UserId {
	return l.userId
}

func (l *Loan) GetBookId() book.BookId {
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

// Repository Interface (Domain Layer)
type ILoanRepository interface {
	Save(loan *Loan) error
	FindById(id LoanId) (*Loan, error)
}

// BorrowingPolicyService - Domain Service for business rule validation
type BorrowingPolicyService struct{}

func NewBorrowingPolicyService() *BorrowingPolicyService {
	return &BorrowingPolicyService{}
}

// CanBorrow validates if a user can borrow a book
func (s *BorrowingPolicyService) CanBorrow(u *user.User, b *book.Book) error {
	// Rule 1: Book must be available
	if !b.IsAvailable() {
		return errors.New("book is not available for borrowing")
	}

	// Rule 2: User must not exceed loan limit
	if !u.CanBorrow() {
		return errors.New("user has reached maximum loan limit")
	}

	return nil
}