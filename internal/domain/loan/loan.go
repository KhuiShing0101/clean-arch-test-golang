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

const (
	LoanPeriodDays = 14
	ExtensionDays  = 7  // Lesson 7: Extend by 7 days
	MaxExtensions  = 2  // Lesson 7: Max 2 extensions per loan
)

type Loan struct {
	id             LoanId
	userId         user.UserId
	bookId         book.BookId
	borrowedAt     time.Time
	dueDate        time.Time
	returnedAt     *time.Time
	extensionCount int // Lesson 7: Track number of extensions
}

// NewLoan - creates new loan (used in BorrowBook)
func NewLoan(userId user.UserId, bookId book.BookId) *Loan {
	now := time.Now()
	dueDate := now.AddDate(0, 0, LoanPeriodDays)

	return &Loan{
		id:             GenerateLoanId(),
		userId:         userId,
		bookId:         bookId,
		borrowedAt:     now,
		dueDate:        dueDate,
		returnedAt:     nil, // Not returned yet
		extensionCount: 0,   // Lesson 7: No extensions yet
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

// ExtendDueDate - extends the loan due date by 7 days (Lesson 7)
func (l *Loan) ExtendDueDate() error {
	// Validation: Loan must be active
	if !l.IsActive() {
		return errors.New("cannot extend: loan is already returned")
	}

	// Validation: Loan must not be overdue
	if l.IsOverdue() {
		return errors.New("cannot extend: loan is overdue")
	}

	// Validation: Cannot exceed max extensions
	if l.extensionCount >= MaxExtensions {
		return errors.New("cannot extend: maximum extension limit reached")
	}

	// Extend due date by 7 days
	l.dueDate = l.dueDate.AddDate(0, 0, ExtensionDays)
	l.extensionCount++

	return nil
}

// IsOverdue - checks if loan is overdue (Lesson 7)
func (l *Loan) IsOverdue() bool {
	if !l.IsActive() {
		return false // Returned loans are not overdue
	}
	return time.Now().After(l.dueDate)
}

// CanExtend - validates if loan can be extended (Lesson 7)
func (l *Loan) CanExtend() bool {
	return l.IsActive() && !l.IsOverdue() && l.extensionCount < MaxExtensions
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

func (l *Loan) GetExtensionCount() int {
	return l.extensionCount
}

// Repository Interface (Domain Layer)
type ILoanRepository interface {
	Save(loan *Loan) error
	FindById(id LoanId) (*Loan, error)
	CreateLoan(userId user.UserId, bookId book.BookId) (*Loan, error)
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