package query

// BookQueryService - Query Service Interface (CQRS Query Side)
//
// Unlike Repository (which works with entities),
// QueryService returns Read Models directly.
type BookQueryService interface {
	// GetBookById retrieves a book by ID with current loan status
	//
	// Returns the book read model if found, nil otherwise
	GetBookById(bookId string) (*BookReadModel, error)

	// ListBooks retrieves paginated books with their loan status (Lesson 6)
	//
	// Parameters:
	//   limit  - Number of items per page
	//   offset - Number of items to skip (calculated from page number)
	//
	// Returns:
	//   books - Array of book read models
	//   total - Total count of all books (for pagination metadata)
	//   error - Any error that occurred
	ListBooks(limit int, offset int) ([]*BookReadModel, int, error)
}