package borrowbook_test

import (
	"testing"
	"time"
	"library-management/internal/application/borrowbook"
	"library-management/internal/domain/user"
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
)

type mockUserRepo struct {
	user *user.User
	err  error
}

func (m *mockUserRepo) FindById(id user.UserId) (*user.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) Save(u *user.User) error { return nil }

type mockBookRepo struct {
	book *book.Book
	err  error
}

func (m *mockBookRepo) FindById(id book.BookId) (*book.Book, error) {
	return m.book, m.err
}

func (m *mockBookRepo) Save(b *book.Book) error { return nil }

type mockLoanRepo struct {
	savedLoan *loan.Loan
}

func (m *mockLoanRepo) Save(l *loan.Loan) error {
	m.savedLoan = l
	return nil
}

func (m *mockLoanRepo) FindById(id loan.LoanId) (*loan.Loan, error) {
	return nil, nil // Not used in tests yet
}

func (m *mockLoanRepo) CreateLoan(userId user.UserId, bookId book.BookId) (*loan.Loan, error) {
	return loan.NewLoan(userId, bookId), nil
}

type mockTxManager struct{}

func (m *mockTxManager) RunInTransaction(fn func() error) error {
	return fn()
}

func TestBorrowBookUseCase_Success(t *testing.T) {
	// Arrange
	u := user.NewUser("user-123", "John Doe", "john@example.com")
	b := book.NewBook("book-456", "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-123", "book-456")

	// Act
	resp, err := useCase.Execute(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if resp.LoanId == "" {
		t.Error("Expected loan ID, got empty string")
	}
}

func TestBorrowBookUseCase_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := &mockUserRepo{user: nil, err: &borrowbook.UserNotFoundError{UserId: "user-999"}}
	bookRepo := &mockBookRepo{}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-999", "book-456")

	// Act
	_, err := useCase.Execute(req)

	// Assert
	if err == nil {
		t.Error("Expected UserNotFoundError")
	}
}

func TestBorrowBookUseCase_BookNotFound(t *testing.T) {
	// Arrange
	u := user.NewUser("user-123", "John Doe", "john@example.com")
	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: nil, err: &borrowbook.BookNotFoundError{BookId: "book-999"}}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-123", "book-999")

	// Act
	_, err := useCase.Execute(req)

	// Assert
	if err == nil {
		t.Error("Expected BookNotFoundError")
	}
}

func TestBorrowBookUseCase_LoanLimitExceeded(t *testing.T) {
	// Arrange
	u := user.NewUser("user-123", "John Doe", "john@example.com")
	// Record 5 loans for user (limit)
	for i := 0; i < 5; i++ {
		u.RecordLoan()
	}

	b := book.NewBook("book-456", "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-123", "book-456")

	// Act
	_, err := useCase.Execute(req)

	// Assert
	if err == nil {
		t.Error("Expected LoanLimitExceededError")
	}
}

func TestBorrowBookUseCase_BookAlreadyBorrowed(t *testing.T) {
	// Arrange
	u := user.NewUser("user-123", "John Doe", "john@example.com")
	b := book.NewBook("book-456", "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")
	b.MarkAsBorrowed() // Mark book as borrowed

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-123", "book-456")

	// Act
	_, err := useCase.Execute(req)

	// Assert
	if err == nil {
		t.Error("Expected BookNotAvailableError")
	}
}

func TestBorrowBookUseCase_CreatesLoanWithCorrectDueDate(t *testing.T) {
	// Arrange
	u := user.NewUser("user-123", "John Doe", "john@example.com")
	b := book.NewBook("book-456", "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")

	userRepo := &mockUserRepo{user: u}
	bookRepo := &mockBookRepo{book: b}
	loanRepo := &mockLoanRepo{}
	policyService := loan.NewBorrowingPolicyService()
	txManager := &mockTxManager{}

	useCase := borrowbook.NewBorrowBookUseCase(
		userRepo, bookRepo, loanRepo, policyService, txManager,
	)

	req, _ := borrowbook.NewBorrowBookRequest("user-123", "book-456")

	// Act
	useCase.Execute(req)

	// Assert - Due date should be 14 days from now
	savedLoan := loanRepo.savedLoan
	if savedLoan == nil {
		t.Fatal("Expected loan to be saved")
	}

	dueDate := savedLoan.GetDueDate()
	now := time.Now()
	expectedDue := now.AddDate(0, 0, 14)

	if dueDate.Format("2006-01-02") != expectedDue.Format("2006-01-02") {
		t.Errorf("Expected due date %v, got %v", expectedDue, dueDate)
	}
}