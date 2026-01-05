package getbook_test

import (
	"testing"
	"library-management/internal/application/getbook"
	"library-management/internal/application/query"
)

// Mock QueryService for testing
type mockQueryService struct {
	bookReadModel *query.BookReadModel
	err           error
}

func (m *mockQueryService) GetBookById(bookId string) (*query.BookReadModel, error) {
	return m.bookReadModel, m.err
}

func (m *mockQueryService) ListBooks(limit int, offset int) ([]*query.BookReadModel, int, error) {
	return nil, 0, nil
}

func TestGetBookUseCase_ReturnsBookWhenAvailable(t *testing.T) {
	// Arrange（準備）
	mockService := &mockQueryService{
		bookReadModel: &query.BookReadModel{
			ID:          "b-12345",
			ISBN:        "978-0-123456-78-9",
			Title:       "Clean Architecture",
			Author:      "Robert C. Martin",
			IsAvailable: true,
			CurrentLoan: nil,
		},
	}

	useCase := getbook.NewGetBookUseCase(mockService)
	request, _ := getbook.NewGetBookRequest("b-12345")

	// Act（実行）
	response := useCase.Execute(request)

	// Assert（検証）
	if !response.Success {
		t.Errorf("Expected success, got failure: %v", *response.ErrorMessage)
	}
	if response.Data == nil {
		t.Error("Expected book data, got nil")
	}
	if response.Data.ID != "b-12345" {
		t.Errorf("Expected book ID 'b-12345', got %s", response.Data.ID)
	}
	if response.Data.CurrentLoan != nil {
		t.Error("Expected no current loan")
	}
}

func TestGetBookUseCase_ReturnsBookWithLoanWhenBorrowed(t *testing.T) {
	// Arrange（準備）
	mockService := &mockQueryService{
		bookReadModel: &query.BookReadModel{
			ID:          "b-12345",
			ISBN:        "978-0-123456-78-9",
			Title:       "Domain-Driven Design",
			Author:      "Eric Evans",
			IsAvailable: false,
			CurrentLoan: &query.LoanReadModel{
				LoanID:       "l-111",
				UserID:       "u-222",
				BorrowedDate: "2024-01-01",
				DueDate:      "2024-01-15",
			},
		},
	}

	useCase := getbook.NewGetBookUseCase(mockService)
	request, _ := getbook.NewGetBookRequest("b-12345")

	// Act（実行）
	response := useCase.Execute(request)

	// Assert（検証）
	if !response.Success {
		t.Errorf("Expected success, got failure")
	}
	if response.Data.CurrentLoan == nil {
		t.Error("Expected current loan, got nil")
	}
}

func TestGetBookUseCase_ReturnsErrorWhenBookNotFound(t *testing.T) {
	// Arrange（準備）
	mockService := &mockQueryService{
		bookReadModel: nil, // 書籍が見つかりません
	}

	useCase := getbook.NewGetBookUseCase(mockService)
	request, _ := getbook.NewGetBookRequest("b-99999")

	// Act（実行）
	response := useCase.Execute(request)

	// Assert（検証）
	if response.Success {
		t.Error("Expected failure, got success")
	}
	if *response.ErrorMessage != "Book not found" {
		t.Errorf("Expected 'Book not found', got %s", *response.ErrorMessage)
	}
}

func TestGetBookRequest_RejectsEmptyBookId(t *testing.T) {
	// Act（実行）
	_, err := getbook.NewGetBookRequest("")

	// Assert（検証）
	if err == nil {
		t.Error("Expected error for empty book ID")
	}
}