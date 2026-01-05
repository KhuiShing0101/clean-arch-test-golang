package borrowbook

import (
	"fmt"

	"library-management/internal/domain/book"
	"library-management/internal/domain/user"
)

// UserNotFoundError - Error when user is not found
type UserNotFoundError struct {
	UserId user.UserId
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found: %s", e.UserId)
}

// BookNotFoundError - Error when book is not found
type BookNotFoundError struct {
	BookId book.BookId
}

func (e *BookNotFoundError) Error() string {
	return fmt.Sprintf("book not found: %s", e.BookId)
}
