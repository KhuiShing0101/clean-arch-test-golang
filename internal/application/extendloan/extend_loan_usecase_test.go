package extendloan_test

import (
	"errors"
	"testing"

	"library-management/internal/application/extendloan"
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
)

// Mock loan repository
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

func (m *mockLoanRepo) CreateLoan(userId user.UserId, bookId book.BookId) (*loan.Loan, error) {
	l := loan.NewLoan(userId, bookId)
	return l, nil
}

// Test 1: Success - First extension (0 → 1)
func TestExtendLoanUseCase_Success_FirstExtension(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")

	testLoan := loan.NewLoan(userId, bookId)
	originalDueDate := testLoan.GetDueDate()

	loanRepo := &mockLoanRepo{loan: testLoan}
	useCase := extendloan.NewExtendLoanUseCase(loanRepo)
	request, _ := extendloan.NewExtendLoanRequest(string(testLoan.GetId()))

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Verify extension count increased
	if testLoan.GetExtensionCount() != 1 {
		t.Errorf("Expected extension count to be 1, got: %d", testLoan.GetExtensionCount())
	}

	// Verify due date extended by 7 days
	expectedDueDate := originalDueDate.AddDate(0, 0, 7)
	if !testLoan.GetDueDate().Equal(expectedDueDate) {
		t.Errorf("Expected due date to be %v, got: %v", expectedDueDate, testLoan.GetDueDate())
	}

	// Verify response data
	if response.ExtensionCount != 1 {
		t.Errorf("Expected response extension count to be 1, got: %d", response.ExtensionCount)
	}

	if response.Message != "Loan extended successfully" {
		t.Errorf("Expected success message, got: %s", response.Message)
	}
}

// Test 2: Success - Second extension (1 → 2)
func TestExtendLoanUseCase_Success_SecondExtension(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")

	testLoan := loan.NewLoan(userId, bookId)

	// First extension
	testLoan.ExtendDueDate()
	dueDateAfterFirst := testLoan.GetDueDate()

	loanRepo := &mockLoanRepo{loan: testLoan}
	useCase := extendloan.NewExtendLoanUseCase(loanRepo)
	request, _ := extendloan.NewExtendLoanRequest(string(testLoan.GetId()))

	// Act - Second extension
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Verify extension count is 2
	if testLoan.GetExtensionCount() != 2 {
		t.Errorf("Expected extension count to be 2, got: %d", testLoan.GetExtensionCount())
	}

	// Verify due date extended by another 7 days
	expectedDueDate := dueDateAfterFirst.AddDate(0, 0, 7)
	if !testLoan.GetDueDate().Equal(expectedDueDate) {
		t.Errorf("Expected due date to be %v, got: %v", expectedDueDate, testLoan.GetDueDate())
	}
}

// Test 3: Error - Loan not found
func TestExtendLoanUseCase_LoanNotFound(t *testing.T) {
	// Arrange
	loanRepo := &mockLoanRepo{loan: nil, err: errors.New("not found")}
	useCase := extendloan.NewExtendLoanUseCase(loanRepo)
	request, _ := extendloan.NewExtendLoanRequest("invalid-loan-id")

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for loan not found")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != extendloan.ErrLoanNotFound {
		t.Errorf("Expected ErrLoanNotFound, got: %v", err)
	}
}

// Test 4: Error - Loan already returned (inactive)
func TestExtendLoanUseCase_LoanAlreadyReturned(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")

	testLoan := loan.NewLoan(userId, bookId)
	testLoan.RecordReturn() // Mark as returned

	loanRepo := &mockLoanRepo{loan: testLoan}
	useCase := extendloan.NewExtendLoanUseCase(loanRepo)
	request, _ := extendloan.NewExtendLoanRequest(string(testLoan.GetId()))

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for returned loan")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	expectedError := "cannot extend: loan is already returned"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got: %v", expectedError, err)
	}
}

// Test 5: Error - Maximum extensions reached
func TestExtendLoanUseCase_MaxExtensionsReached(t *testing.T) {
	// Arrange
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")

	testLoan := loan.NewLoan(userId, bookId)

	// Use up both extensions
	testLoan.ExtendDueDate() // First extension
	testLoan.ExtendDueDate() // Second extension

	if testLoan.GetExtensionCount() != 2 {
		t.Fatalf("Setup failed: expected 2 extensions, got: %d", testLoan.GetExtensionCount())
	}

	loanRepo := &mockLoanRepo{loan: testLoan}
	useCase := extendloan.NewExtendLoanUseCase(loanRepo)
	request, _ := extendloan.NewExtendLoanRequest(string(testLoan.GetId()))

	// Act - Try third extension
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for max extensions reached")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	expectedError := "cannot extend: maximum extension limit reached"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got: %v", expectedError, err)
	}

	// Verify extension count didn't increase
	if testLoan.GetExtensionCount() != 2 {
		t.Errorf("Expected extension count to remain 2, got: %d", testLoan.GetExtensionCount())
	}
}

// Test 6: Verify CanExtend() validation logic
func TestLoanCanExtend(t *testing.T) {
	userId := user.UserId("user-123")
	bookId := book.BookId("book-456")

	// Test case 1: Fresh loan can be extended
	testLoan := loan.NewLoan(userId, bookId)
	if !testLoan.CanExtend() {
		t.Error("Expected fresh loan to be extendable")
	}

	// Test case 2: After 1 extension, can still extend
	testLoan.ExtendDueDate()
	if !testLoan.CanExtend() {
		t.Error("Expected loan with 1 extension to be extendable")
	}

	// Test case 3: After 2 extensions, cannot extend
	testLoan.ExtendDueDate()
	if testLoan.CanExtend() {
		t.Error("Expected loan with 2 extensions to not be extendable")
	}

	// Test case 4: Returned loan cannot be extended
	returnedLoan := loan.NewLoan(userId, bookId)
	returnedLoan.RecordReturn()
	if returnedLoan.CanExtend() {
		t.Error("Expected returned loan to not be extendable")
	}
}

// Test 7: Request validation
func TestExtendLoanRequest_Validation(t *testing.T) {
	// Test case 1: Valid request
	validRequest, err := extendloan.NewExtendLoanRequest("valid-loan-id")
	if err != nil {
		t.Errorf("Expected valid request to succeed, got error: %v", err)
	}
	if validRequest == nil {
		t.Error("Expected valid request object")
	}

	// Test case 2: Empty loan ID
	invalidRequest, err := extendloan.NewExtendLoanRequest("")
	if err == nil {
		t.Error("Expected error for empty loan ID")
	}
	if invalidRequest != nil {
		t.Error("Expected nil request for invalid input")
	}

	expectedError := "loan ID is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got: %v", expectedError, err)
	}
}
