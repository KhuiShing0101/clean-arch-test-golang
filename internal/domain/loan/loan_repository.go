package loan

import (
	"library-management/internal/domain/book"
	"library-management/internal/domain/user"
)

type LoanRepository interface {
	Save(loan *Loan) error
	FindById(id *LoanId) (*Loan, error)
	FindByUserId(userId *user.UserId) ([]*Loan, error)
	FindByBookId(bookId *book.BookId) ([]*Loan, error)
	FindActiveLoans() ([]*Loan, error)
	FindOverdueLoans() ([]*Loan, error)
	Delete(id *LoanId) error
}