package listbooks

import (
	"library-management/internal/application/getbook"
)

// ListBooksResponse - DTO for list books result
type ListBooksResponse struct {
	Books      []getbook.BookDTO
	Pagination PaginationMetadata
}

// PaginationMetadata - Metadata about the paginated result
type PaginationMetadata struct {
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

// NewListBooksResponse creates a new response
func NewListBooksResponse(
	books []getbook.BookDTO,
	page int,
	limit int,
	total int,
) *ListBooksResponse {
	totalPages := (total + limit - 1) / limit // Ceiling division
	if totalPages == 0 {
		totalPages = 1
	}

	return &ListBooksResponse{
		Books: books,
		Pagination: PaginationMetadata{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
