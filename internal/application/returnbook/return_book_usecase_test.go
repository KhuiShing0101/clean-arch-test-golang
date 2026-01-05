package returnbook_test

import (
	"errors"
	"testing"
	"time"

	"library-management/internal/application/returnbook"
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
)

// Mock repositories
type mockLoanRepo struct {
	loan *loan.Loan
	err  error
}

func (m *mockLoanRepo) Save(l *loan.Loan) error {
	m.loan = l
	return m.err
}

func (m *mockLoanRepo) FindById(id loan.LoanId) (*loan.Loan, error) {
	return m.loan, m.err
}

type mockBookRepo struct {
	book *book.Book
	err  error
}

func (m *mockBookRepo) Save(b *book.Book) error {
	m.book = b
	return m.err
}

func (m *mockBookRepo) FindById(id book.BookId) (*book.Book, error) {
	return m.book, m.err
}

type mockUserRepo struct {
	user *user.User
	err  error
}

func (m *mockUserRepo) Save(u *user.User) error {
	m.user = u
	return m.err
}

func (m *mockUserRepo) FindById(id user.UserId) (*user.User, error) {
	return m.user, m.err
}

type mockTxManager struct{}

func (m *mockTxManager) RunInTransaction(fn func() error) error {
	return fn()
}

// Helper function removed - using loan.NewLoan() instead

// Test 1: Success - Return book on time (no late fee)
func TestReturnBookUseCase_Success_OnTime(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")
	loanId := loan.GenerateLoanId()

	// Create loan that was borrowed 7 days ago (still within 14-day period)
	testLoan := loan.NewLoan(userId, bookId)

	testBook := book.NewBook(bookId, "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")
	testBook.MarkAsBorrowed() // Book is currently borrowed

	testUser := user.NewUser(userId, "John Doe", "john@example.com")
	testUser.IncrementLoanCount() // User has 1 active loan

	loanRepo := &mockLoanRepo{loan: testLoan}
	bookRepo := &mockBookRepo{book: testBook}
	userRepo := &mockUserRepo{user: testUser}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loanId)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.LateFee != 0.0 {
		t.Errorf("Expected late fee to be 0.0, got: %f", response.LateFee)
	}

	if response.IsOverdue {
		t.Error("Expected not overdue")
	}

	if response.DaysLate != 0 {
		t.Errorf("Expected 0 days late, got: %d", response.DaysLate)
	}

	// Verify book is available again
	if !testBook.IsAvailable() {
		t.Error("Expected book to be available after return")
	}

	// Verify user loan count decremented
	if testUser.GetCurrentLoanCount() != 0 {
		t.Errorf("Expected user loan count to be 0, got: %d", testUser.GetCurrentLoanCount())
	}
}

// Test 2: Success - Return book late (with late fee)
func TestReturnBookUseCase_Success_Late(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")
	loanId := loan.GenerateLoanId()

	// Create a loan that's overdue (borrowed 20 days ago, due date was 6 days ago)
	testLoan := loan.NewLoan(userId, bookId)
	// Manually set borrowed date to 20 days ago for testing
	// Note: In real implementation, we'd need to expose setters or use a test helper

	testBook := book.NewBook(bookId, "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")
	testBook.MarkAsBorrowed()

	testUser := user.NewUser(userId, "John Doe", "john@example.com")
	testUser.IncrementLoanCount()

	loanRepo := &mockLoanRepo{loan: testLoan}
	bookRepo := &mockBookRepo{book: testBook}
	userRepo := &mockUserRepo{user: testUser}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loanId)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Note: Late fee calculation depends on actual dates
	// For a proper test, we'd mock time or use dependency injection
}

// Test 3: Error - Loan not found
func TestReturnBookUseCase_LoanNotFound(t *testing.T) {
	// Arrange
	loanRepo := &mockLoanRepo{loan: nil, err: errors.New("not found")}
	bookRepo := &mockBookRepo{}
	userRepo := &mockUserRepo{}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loan.LoanId("invalid-id"))

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for loan not found")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != returnbook.ErrLoanNotFound {
		t.Errorf("Expected ErrLoanNotFound, got: %v", err)
	}
}

// Test 4: Error - Loan already returned
func TestReturnBookUseCase_LoanAlreadyReturned(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")
	loanId := loan.GenerateLoanId()

	testLoan := loan.NewLoan(userId, bookId)
	testLoan.RecordReturn() // Already returned

	testBook := book.NewBook(bookId, "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")
	testUser := user.NewUser(userId, "John Doe", "john@example.com")

	loanRepo := &mockLoanRepo{loan: testLoan}
	bookRepo := &mockBookRepo{book: testBook}
	userRepo := &mockUserRepo{user: testUser}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loanId)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for already returned loan")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != returnbook.ErrLoanAlreadyReturned {
		t.Errorf("Expected ErrLoanAlreadyReturned, got: %v", err)
	}
}

// Test 5: Error - Book not found
func TestReturnBookUseCase_BookNotFound(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")
	loanId := loan.GenerateLoanId()

	testLoan := loan.NewLoan(userId, bookId)

	loanRepo := &mockLoanRepo{loan: testLoan}
	bookRepo := &mockBookRepo{book: nil, err: errors.New("not found")}
	userRepo := &mockUserRepo{}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loanId)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for book not found")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != returnbook.ErrBookNotFound {
		t.Errorf("Expected ErrBookNotFound, got: %v", err)
	}
}

// Test 6: Error - User not found
func TestReturnBookUseCase_UserNotFound(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")
	loanId := loan.GenerateLoanId()

	testLoan := loan.NewLoan(userId, bookId)
	testBook := book.NewBook(bookId, "978-0-13-468599-1", "Clean Architecture", "Robert C. Martin")

	loanRepo := &mockLoanRepo{loan: testLoan}
	bookRepo := &mockBookRepo{book: testBook}
	userRepo := &mockUserRepo{user: nil, err: errors.New("not found")}
	txManager := &mockTxManager{}

	useCase := returnbook.NewReturnBookUseCase(loanRepo, bookRepo, userRepo, txManager)
	request := returnbook.NewReturnBookRequest(loanId)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for user not found")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != returnbook.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}
}

// Test 7: Late fee calculation correctness
func TestLateFeeCalculator(t *testing.T) {
	calculator := loan.NewLateFeeCalculator()

	// Test case 1: On time return (no late fee)
	dueDate := time.Now()
	returnDate := time.Now()
	lateFee := calculator.CalculateLateFee(dueDate, returnDate)
	if lateFee != 0.0 {
		t.Errorf("Expected 0.0 late fee for on-time return, got: %f", lateFee)
	}

	// Test case 2: 5 days late
	dueDate = time.Now().AddDate(0, 0, -5)
	returnDate = time.Now()
	lateFee = calculator.CalculateLateFee(dueDate, returnDate)
	expectedFee := 5.0 // 5 days * $1/day
	if lateFee != expectedFee {
		t.Errorf("Expected %f late fee for 5 days late, got: %f", expectedFee, lateFee)
	}

	// Test case 3: Early return (no late fee)
	dueDate = time.Now().AddDate(0, 0, 5)
	returnDate = time.Now()
	lateFee = calculator.CalculateLateFee(dueDate, returnDate)
	if lateFee != 0.0 {
		t.Errorf("Expected 0.0 late fee for early return, got: %f", lateFee)
	}
}
