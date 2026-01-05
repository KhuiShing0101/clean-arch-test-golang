package returnbook

import (
	"errors"
	"time"

	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/shared"
	"library-management/internal/domain/user"
)

// Custom errors
var (
	ErrLoanNotFound      = errors.New("loan not found")
	ErrLoanAlreadyReturned = errors.New("loan is already returned")
	ErrBookNotFound      = errors.New("book not found")
	ErrUserNotFound      = errors.New("user not found")
)

// ReturnBookUseCase - Application layer use case for returning a book
type ReturnBookUseCase struct {
	loanRepo         loan.ILoanRepository
	bookRepo         book.IBookRepository
	userRepo         user.IUserRepository
	txManager        shared.TransactionManager
	lateFeeCalculator loan.ILateFeeCalculator
}

// NewReturnBookUseCase creates a new use case instance
// The lateFeeCalculator is injected as a dependency (Dependency Inversion Principle)
func NewReturnBookUseCase(
	loanRepo loan.ILoanRepository,
	bookRepo book.IBookRepository,
	userRepo user.IUserRepository,
	txManager shared.TransactionManager,
	lateFeeCalculator loan.ILateFeeCalculator,
) *ReturnBookUseCase {
	return &ReturnBookUseCase{
		loanRepo:         loanRepo,
		bookRepo:         bookRepo,
		userRepo:         userRepo,
		txManager:        txManager,
		lateFeeCalculator: lateFeeCalculator,
	}
}

// Execute processes the return book request
func (uc *ReturnBookUseCase) Execute(req *ReturnBookRequest) (*ReturnBookResponse, error) {
	var response *ReturnBookResponse
	var execError error

	// Run in transaction
	txErr := uc.txManager.RunInTransaction(func() error {
		// 1. Find the loan
		loanEntity, err := uc.loanRepo.FindById(req.LoanId)
		if err != nil || loanEntity == nil {
			execError = ErrLoanNotFound
			return execError
		}

		// 2. Verify loan is active (not already returned)
		if !loanEntity.IsActive() {
			execError = ErrLoanAlreadyReturned
			return execError
		}

		// 3. Find the book
		bookEntity, err := uc.bookRepo.FindById(loanEntity.GetBookId())
		if err != nil || bookEntity == nil {
			execError = ErrBookNotFound
			return execError
		}

		// 4. Find the user
		userEntity, err := uc.userRepo.FindById(loanEntity.GetUserId())
		if err != nil || userEntity == nil {
			execError = ErrUserNotFound
			return execError
		}

		// 5. Record return on loan entity
		now := time.Now()
		if err := loanEntity.RecordReturn(); err != nil {
			execError = err
			return execError
		}

		// 6. Calculate late fee using domain service
		dueDate := loanEntity.GetDueDate()
		lateFee := uc.lateFeeCalculator.CalculateLateFee(dueDate, now)
		daysLate := uc.lateFeeCalculator.GetDaysLate(dueDate, now)
		isOverdue := uc.lateFeeCalculator.IsOverdue(dueDate, now)

		// 7. Update book status (mark as available again)
		bookEntity.MarkAsAvailable()

		// 8. Update user (decrement current loan count)
		userEntity.RecordReturn()

		// 9. Save all changes
		if err := uc.loanRepo.Save(loanEntity); err != nil {
			execError = err
			return execError
		}

		if err := uc.bookRepo.Save(bookEntity); err != nil {
			execError = err
			return execError
		}

		if err := uc.userRepo.Save(userEntity); err != nil {
			execError = err
			return execError
		}

		// 10. Create response
		response = NewReturnBookResponse(
			loanEntity.GetId(),
			loanEntity.GetBookId(),
			loanEntity.GetUserId(),
			loanEntity.GetBorrowedAt(),
			loanEntity.GetDueDate(),
			*loanEntity.GetReturnedAt(),
			daysLate,
			lateFee,
			isOverdue,
		)

		return nil
	})

	if txErr != nil {
		if execError != nil {
			return nil, execError
		}
		return nil, txErr
	}

	return response, nil
}
