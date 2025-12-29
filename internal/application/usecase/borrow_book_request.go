package usecase

import "errors"

// BorrowBookRequest - Input for BorrowBook use case
// Simple data carrier - no business logic!
type BorrowBookRequest struct {
	// WHY separate struct? Keeps use case signature clean
	// Single parameter instead of multiple primitive arguments
	UserId string
	BookId string
}

// NewBorrowBookRequest - バリデーション付きで新しいBorrowBookRequestを作成
func NewBorrowBookRequest(userId, bookId string) (*BorrowBookRequest, error) {
	// 基本的なバリデーション - 空でないかチェック
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