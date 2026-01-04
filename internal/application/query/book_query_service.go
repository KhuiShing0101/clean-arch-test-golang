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

	// ListBooks retrieves all books with their loan status (for Lesson 6)
	//
	// Returns array of book read models
	ListBooks() ([]*BookReadModel, error)
}