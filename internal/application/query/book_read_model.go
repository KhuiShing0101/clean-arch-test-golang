package query

// BookReadModel - Read Model for book data (CQRS Query Side)
//
// This is NOT a domain entity! It's a simple struct
// optimized for read operations.
type BookReadModel struct {
	ID           string          // Book ID
	ISBN         string          // Book ISBN
	Title        string          // Book title
	Author       string          // Book author
	IsAvailable  bool            // Derived from loans!
	CurrentLoan  *LoanReadModel  // Current loan info if borrowed
}

// LoanReadModel - Read Model for loan data
type LoanReadModel struct {
	LoanID       string  // Loan ID
	UserID       string  // User ID who borrowed
	BorrowedDate string  // When book was borrowed
	DueDate      string  // When book is due
}