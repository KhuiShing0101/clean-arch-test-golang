package listbooks_test

import (
	"errors"
	"fmt"
	"testing"

	"library-management/internal/application/getbook"
	"library-management/internal/application/listbooks"
)

// Mock BookQueryService
type mockBookQueryService struct {
	books []getbook.BookDTO
	total int
	err   error
}

func (m *mockBookQueryService) GetBook(id string) (*getbook.BookDTO, error) {
	// Not used in ListBooks
	return nil, nil
}

func (m *mockBookQueryService) ListBooks(limit int, offset int) ([]getbook.BookDTO, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	// Simulate pagination
	start := offset
	end := offset + limit
	if start > len(m.books) {
		return []getbook.BookDTO{}, m.total, nil
	}
	if end > len(m.books) {
		end = len(m.books)
	}

	return m.books[start:end], m.total, nil
}

// Test 1: Success - First page with default limit
func TestListBooksUseCase_Success_FirstPage(t *testing.T) {
	// Arrange
	testBooks := []getbook.BookDTO{
		{BookId: "b-1", ISBN: "isbn-1", Title: "Book 1", Author: "Author 1", IsAvailable: true},
		{BookId: "b-2", ISBN: "isbn-2", Title: "Book 2", Author: "Author 2", IsAvailable: false},
		{BookId: "b-3", ISBN: "isbn-3", Title: "Book 3", Author: "Author 3", IsAvailable: true},
	}

	queryService := &mockBookQueryService{
		books: testBooks,
		total: len(testBooks),
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(1, 10)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(response.Books) != 3 {
		t.Errorf("Expected 3 books, got: %d", len(response.Books))
	}

	if response.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got: %d", response.Pagination.Page)
	}

	if response.Pagination.Total != 3 {
		t.Errorf("Expected total 3, got: %d", response.Pagination.Total)
	}
}

// Test 2: Success - Second page
func TestListBooksUseCase_Success_SecondPage(t *testing.T) {
	// Arrange - Create 15 books
	testBooks := make([]getbook.BookDTO, 15)
	for i := 0; i < 15; i++ {
		testBooks[i] = getbook.BookDTO{
			BookId:      fmt.Sprintf("b-%d", i+1),
			ISBN:        fmt.Sprintf("isbn-%d", i+1),
			Title:       fmt.Sprintf("Book %d", i+1),
			Author:      fmt.Sprintf("Author %d", i+1),
			IsAvailable: true,
		}
	}

	queryService := &mockBookQueryService{
		books: testBooks,
		total: len(testBooks),
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(2, 10)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Second page should have 5 books (15 total - 10 on first page)
	if len(response.Books) != 5 {
		t.Errorf("Expected 5 books on page 2, got: %d", len(response.Books))
	}

	if response.Pagination.Page != 2 {
		t.Errorf("Expected page 2, got: %d", response.Pagination.Page)
	}

	if response.Pagination.TotalPages != 2 {
		t.Errorf("Expected 2 total pages, got: %d", response.Pagination.TotalPages)
	}
}

// Test 3: Success - Custom limit
func TestListBooksUseCase_Success_CustomLimit(t *testing.T) {
	// Arrange
	testBooks := make([]getbook.BookDTO, 10)
	for i := 0; i < 10; i++ {
		testBooks[i] = getbook.BookDTO{
			BookId:      fmt.Sprintf("b-%d", i+1),
			ISBN:        fmt.Sprintf("isbn-%d", i+1),
			Title:       fmt.Sprintf("Book %d", i+1),
			Author:      fmt.Sprintf("Author %d", i+1),
			IsAvailable: true,
		}
	}

	queryService := &mockBookQueryService{
		books: testBooks,
		total: len(testBooks),
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(1, 5)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(response.Books) != 5 {
		t.Errorf("Expected 5 books with limit=5, got: %d", len(response.Books))
	}

	if response.Pagination.Limit != 5 {
		t.Errorf("Expected limit 5, got: %d", response.Pagination.Limit)
	}

	if response.Pagination.TotalPages != 2 {
		t.Errorf("Expected 2 total pages (10 books / 5 per page), got: %d", response.Pagination.TotalPages)
	}
}

// Test 4: Error - Invalid page (zero)
func TestListBooksUseCase_Error_InvalidPageZero(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{}
	useCase := listbooks.NewListBooksUseCase(queryService)
	request := &listbooks.ListBooksRequest{Page: 0, Limit: 10}

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid page")
	}

	if response != nil {
		t.Error("Expected nil response")
	}

	if err != listbooks.ErrInvalidPage {
		t.Errorf("Expected ErrInvalidPage, got: %v", err)
	}
}

// Test 5: Error - Invalid page (negative)
func TestListBooksUseCase_Error_InvalidPageNegative(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{}
	useCase := listbooks.NewListBooksUseCase(queryService)
	request := &listbooks.ListBooksRequest{Page: -1, Limit: 10}

	// Act
	_, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid page")
	}

	if err != listbooks.ErrInvalidPage {
		t.Errorf("Expected ErrInvalidPage, got: %v", err)
	}
}

// Test 6: Error - Invalid limit (zero)
func TestListBooksUseCase_Error_InvalidLimitZero(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{}
	useCase := listbooks.NewListBooksUseCase(queryService)
	request := &listbooks.ListBooksRequest{Page: 1, Limit: 0}

	// Act
	_, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid limit")
	}

	if err != listbooks.ErrInvalidLimit {
		t.Errorf("Expected ErrInvalidLimit, got: %v", err)
	}
}

// Test 7: Error - Invalid limit (negative)
func TestListBooksUseCase_Error_InvalidLimitNegative(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{}
	useCase := listbooks.NewListBooksUseCase(queryService)
	request := &listbooks.ListBooksRequest{Page: 1, Limit: -10}

	// Act
	_, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid limit")
	}

	if err != listbooks.ErrInvalidLimit {
		t.Errorf("Expected ErrInvalidLimit, got: %v", err)
	}
}

// Test 8: Success - Out of range page (returns empty list)
func TestListBooksUseCase_Success_OutOfRangePage(t *testing.T) {
	// Arrange
	testBooks := []getbook.BookDTO{
		{BookId: "b-1", ISBN: "isbn-1", Title: "Book 1", Author: "Author 1", IsAvailable: true},
	}

	queryService := &mockBookQueryService{
		books: testBooks,
		total: len(testBooks),
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(10, 10) // Page 10 when only 1 book exists

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for out of range page, got: %v", err)
	}

	if len(response.Books) != 0 {
		t.Errorf("Expected 0 books for out of range page, got: %d", len(response.Books))
	}

	if response.Pagination.Total != 1 {
		t.Errorf("Expected total 1, got: %d", response.Pagination.Total)
	}
}

// Test 9: Success - Empty book list
func TestListBooksUseCase_Success_EmptyList(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{
		books: []getbook.BookDTO{},
		total: 0,
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(1, 10)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for empty list, got: %v", err)
	}

	if len(response.Books) != 0 {
		t.Errorf("Expected 0 books, got: %d", len(response.Books))
	}

	if response.Pagination.Total != 0 {
		t.Errorf("Expected total 0, got: %d", response.Pagination.Total)
	}
}

// Test 10: Error - Query service error
func TestListBooksUseCase_Error_QueryServiceFailure(t *testing.T) {
	// Arrange
	queryService := &mockBookQueryService{
		err: errors.New("database connection failed"),
	}

	useCase := listbooks.NewListBooksUseCase(queryService)
	request := listbooks.NewListBooksRequest(1, 10)

	// Act
	response, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error from query service")
	}

	if response != nil {
		t.Error("Expected nil response")
	}
}
