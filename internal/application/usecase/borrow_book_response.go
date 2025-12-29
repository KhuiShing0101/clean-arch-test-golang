package usecase

import "time"

// BorrowBookResponse - Output for BorrowBook use case
// Contains all information the caller needs to know
type BorrowBookResponse struct {
	Success      bool
	LoanId       *string
	DueDate      *time.Time
	ErrorMessage *string
}

// NewSuccessResponse - Factory function for success case
func NewSuccessResponse(loanId string, dueDate time.Time) *BorrowBookResponse {
	return &BorrowBookResponse{
		Success:      true,
		LoanId:       &loanId,
		DueDate:      &dueDate,
		ErrorMessage: nil,
	}
}

// NewFailureResponse - Factory function for failure case
func NewFailureResponse(errorMessage string) *BorrowBookResponse {
	return &BorrowBookResponse{
		Success:      false,
		LoanId:       nil,
		DueDate:      nil,
		ErrorMessage: &errorMessage,
	}
}

// IsSuccess - レスポンスが成功操作を表すかチェック
func (r *BorrowBookResponse) IsSuccess() bool {
	return r.Success
}

// GetLoanId - 利用可能な場合、貸出IDを返す
func (r *BorrowBookResponse) GetLoanId() *string {
	return r.LoanId
}

// GetDueDate - 利用可能な場合、返却期限を返す
func (r *BorrowBookResponse) GetDueDate() *time.Time {
	return r.DueDate
}

// GetErrorMessage - 利用可能な場合、エラーメッセージを返す
func (r *BorrowBookResponse) GetErrorMessage() *string {
	return r.ErrorMessage
}