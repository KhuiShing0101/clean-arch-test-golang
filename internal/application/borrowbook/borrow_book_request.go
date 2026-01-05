package borrowbook

import (
	"errors"

	"library-management/internal/domain/book"
	"library-management/internal/domain/user"
)

type BorrowBookRequest struct {
	UserId user.UserId
	BookId book.BookId
}

func NewBorrowBookRequest(userId user.UserId, bookId book.BookId) (*BorrowBookRequest, error) {
	if userId == "" {
		return nil, errors.New("userId cannot be empty")
	}
	if bookId == "" {
		return nil, errors.New("bookId cannot be empty")
	}
	return &BorrowBookRequest{
		UserId: userId,
		BookId: bookId,
	}, nil
}