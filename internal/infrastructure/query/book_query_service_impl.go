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

func (s *BookQueryServiceImpl) ListBooks(limit int, offset int) ([]*appQuery.BookReadModel, int, error) {
	// Query paginated books with JOIN to get availability
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
		ORDER BY b.id
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	books := []*appQuery.BookReadModel{}
	for rows.Next() {
		var (
			id, isbn, title, author string
			loanId, userId          sql.NullString
			borrowedAt, dueDate     sql.NullString
		)

		err := rows.Scan(
			&id, &isbn, &title, &author,
			&loanId, &userId, &borrowedAt, &dueDate,
		)
		if err != nil {
			return nil, 0, err
		}

		var currentLoan *appQuery.LoanReadModel
		if loanId.Valid {
			currentLoan = &appQuery.LoanReadModel{
				LoanID:       loanId.String,
				UserID:       userId.String,
				BorrowedDate: borrowedAt.String,
				DueDate:      dueDate.String,
			}
		}

		books = append(books, &appQuery.BookReadModel{
			ID:          id,
			ISBN:        isbn,
			Title:       title,
			Author:      author,
			IsAvailable: currentLoan == nil,
			CurrentLoan: currentLoan,
		})
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM books`
	err = s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}