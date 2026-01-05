package borrowbook_test

import (
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
)

// Mock Loan Repository with CreateLoan method
type mockLoanRepository struct {
	loans map[loan.LoanId]*loan.Loan
}

func newMockLoanRepository() *mockLoanRepository {
	return &mockLoanRepository{
		loans: make(map[loan.LoanId]*loan.Loan),
	}
}

func (m *mockLoanRepository) Save(l *loan.Loan) error {
	m.loans[l.GetId()] = l
	return nil
}

func (m *mockLoanRepository) FindById(id loan.LoanId) (*loan.Loan, error) {
	return m.loans[id], nil
}

func (m *mockLoanRepository) CreateLoan(userId user.UserId, bookId book.BookId) (*loan.Loan, error) {
	// Create the loan entity using domain factory
	l := loan.NewLoan(userId, bookId)
	return l, nil
}
