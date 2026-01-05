package listbooks

import (
	"errors"

	"library-management/internal/application/getbook"
)

// Custom errors
var (
	ErrInvalidPage  = errors.New("invalid page number")
	ErrInvalidLimit = errors.New("invalid limit")
)

// ListBooksUseCase - Application layer use case for listing books with pagination
type ListBooksUseCase struct {
	queryService getbook.IBookQueryService
}

// NewListBooksUseCase creates a new use case instance
func NewListBooksUseCase(queryService getbook.IBookQueryService) *ListBooksUseCase {
	return &ListBooksUseCase{
		queryService: queryService,
	}
}

// Execute processes the list books request
func (uc *ListBooksUseCase) Execute(req *ListBooksRequest) (*ListBooksResponse, error) {
	// Validate request
	if req.Page <= 0 {
		return nil, ErrInvalidPage
	}
	if req.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	// Calculate offset
	offset := req.GetOffset()

	// Get paginated books from query service
	books, total, err := uc.queryService.ListBooks(req.Limit, offset)
	if err != nil {
		return nil, err
	}

	// Create response with pagination metadata
	response := NewListBooksResponse(books, req.Page, req.Limit, total)

	return response, nil
}
