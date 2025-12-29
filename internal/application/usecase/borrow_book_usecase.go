package usecase

import (
	"fmt"
	"time"
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/user"
)

type BorrowBookUseCase struct {
	// No TransactionManager needed - BorrowBook has single write!
	userRepository user.UserRepository
	bookRepository book.BookRepository
	loanRepository loan.LoanRepository
	eligibilityService *loan.LoanEligibilityService
}

func (uc *BorrowBookUseCase) Execute(request *BorrowBookRequest) *BorrowBookResponse {
	// Load entities - CONVERT STRING IDs to VALUE OBJECTS
	userId, err := user.NewUserId(request.UserId)
	if err != nil {
		return NewFailureResponse("Invalid user ID: " + err.Error())
	}

	userEntity, err := uc.userRepository.FindById(userId)
	if err != nil {
		return NewFailureResponse(err.Error())
	}

	bookId, err := book.NewBookId(request.BookId)
	if err != nil {
		return NewFailureResponse("Invalid book ID: " + err.Error())
	}

	bookEntity, err := uc.bookRepository.FindById(bookId)
	if err != nil {
		return NewFailureResponse(err.Error())
	}

	// Check eligibility
	if !uc.eligibilityService.CanBorrow(userEntity, bookEntity) {
		reason := uc.eligibilityService.GetIneligibilityReason(userEntity, bookEntity)
		if reason != nil {
			return NewFailureResponse(*reason)
		}
		return NewFailureResponse("User cannot borrow this book")
	}

	// Create loan - GENERATE LOAN ID
	loanId, err := loan.NewLoanId(fmt.Sprintf("LOAN-%d", time.Now().Unix()))
	if err != nil {
		return NewFailureResponse("Failed to generate loan ID: " + err.Error())
	}

	loanEntity := loan.NewLoan(
		loanId,
		userEntity.Id(),      // User has Id() method
		bookEntity.GetId(),   // Book has GetId() method
		time.Now(),
		nil, // returnedAt is nil for new loans
	)

	// Save loan (single write - no transaction needed for BorrowBook!)
	if err := uc.loanRepository.Save(loanEntity); err != nil {
		return NewFailureResponse(err.Error())
	}
	// Book availability is automatically derived from loans table

	return NewSuccessResponse(
		loanEntity.GetId().GetValue(),
		loanEntity.GetDueDate(),
	)
}
