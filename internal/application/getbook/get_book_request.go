package getbook

import "errors"

type GetBookRequest struct {
	BookId string
}

func NewGetBookRequest(bookId string) (*GetBookRequest, error) {
	if bookId == "" {
		return nil, errors.New("bookId cannot be empty")
	}
	return &GetBookRequest{BookId: bookId}, nil
}