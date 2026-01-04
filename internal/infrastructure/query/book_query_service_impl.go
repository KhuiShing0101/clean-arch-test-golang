package query

import (
	"database/sql"
	appQuery "library-management/internal/application/query"
)

// BookQueryServiceImpl implements BookQueryService interface
type BookQueryServiceImpl struct {
	db *sql.DB
}

func NewBookQueryService(db *sql.DB) *BookQueryServiceImpl {
	return &BookQueryServiceImpl{db: db}
}

func (s *BookQueryServiceImpl) GetBookById(bookId string) (*appQuery.BookReadModel, error) {
	// Single JOIN query - efficient!
	query := `
		SELECT
			b.id,
			b.isbn,
			b.title,
			b.author,
			l.id as loan_id,
			l.user_id,
			l.borrowed_at,
			l.due_date
		FROM books b
		LEFT JOIN loans l ON b.id = l.book_id
			AND l.returned_at IS NULL
		WHERE b.id = ?
	`

	var (
		id, isbn, title, author string
		loanId, userId          sql.NullString
		borrowedAt, dueDate     sql.NullString
	)

	err := s.db.QueryRow(query, bookId).Scan(
		&id, &isbn, &title, &author,
		&loanId, &userId, &borrowedAt, &dueDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil  // Book not found
	}
	if err != nil {
		return nil, err
	}

	// Build Read Model directly from query result
	var currentLoan *appQuery.LoanReadModel
	if loanId.Valid {
		currentLoan = &appQuery.LoanReadModel{
			LoanID:       loanId.String,
			UserID:       userId.String,
			BorrowedDate: borrowedAt.String,
			DueDate:      dueDate.String,
		}
	}

	return &appQuery.BookReadModel{
		ID:          id,
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		IsAvailable: currentLoan == nil,  // Single Source of Truth!
		CurrentLoan: currentLoan,
	}, nil
}

func (s *BookQueryServiceImpl) ListBooks() ([]*appQuery.BookReadModel, error) {
	// Will implement in Lesson 6
	return []*appQuery.BookReadModel{}, nil
}