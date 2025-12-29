package borrowbook

import (
	"context"
	"time"

	"myapp/internal/domain/entity"
)

// BorrowBookUseCase orchestrates the book borrowing process
type BorrowBookUseCase struct {
	bookRepository BookRepositoryInterface
	userRepository UserRepositoryInterface
	loanRepository LoanRepositoryInterface

	// 3 domain services (all need repo access)
	loanLimitChecker        *LoanLimitChecker
	overdueLoanChecker      *OverdueLoanChecker
	bookAvailabilityChecker *BookAvailabilityChecker
}

// NewBorrowBookUseCase creates a new BorrowBookUseCase with all dependencies
func NewBorrowBookUseCase(
	bookRepo BookRepositoryInterface,
	userRepo UserRepositoryInterface,
	loanRepo LoanRepositoryInterface,
	limitChecker *LoanLimitChecker,
	overdueChecker *OverdueLoanChecker,
	availabilityChecker *BookAvailabilityChecker,
) *BorrowBookUseCase {
	return &BorrowBookUseCase{
		bookRepository:          bookRepo,
		userRepository:          userRepo,
		loanRepository:          loanRepo,
		loanLimitChecker:        limitChecker,
		overdueLoanChecker:      overdueChecker,
		bookAvailabilityChecker: availabilityChecker,
	}
}

// Execute runs the borrow book use case
func (uc *BorrowBookUseCase) Execute(ctx context.Context, req *BorrowBookRequest) (*BorrowBookResponse, error) {
	// STEP 1: Load entities from repositories
	user, err := uc.userRepository.FindByID(ctx, req.UserID)
	if err != nil {
		return NewBorrowBookResponseFailure("user not found"), nil
	}

	book, err := uc.bookRepository.FindByID(ctx, req.BookID)
	if err != nil {
		return NewBorrowBookResponseFailure("book not found"), nil
	}

	// STEP 2: Check user eligibility (ENTITY METHOD - no repo needed!)
	if !user.CanBorrow() {
		return NewBorrowBookResponseFailure("user cannot borrow: account not active"), nil
	}

	// STEP 3: Check loan limit (DOMAIN SERVICE)
	activeLoans, _ := uc.loanRepository.FindActiveLoansByUserID(ctx, user.GetID())
	if err := uc.loanLimitChecker.Check(user, activeLoans); err != nil {
		return NewBorrowBookResponseFailure(err.Error()), nil
	}

	// STEP 4: Check for overdue loans (DOMAIN SERVICE)
	overdueLoans, _ := uc.loanRepository.FindOverdueLoansByUserID(ctx, user.GetID())
	if err := uc.overdueLoanChecker.Check(overdueLoans); err != nil {
		return NewBorrowBookResponseFailure(err.Error()), nil
	}

	// STEP 5: Check book availability (DOMAIN SERVICE - Single Source of Truth!)
	if err := uc.bookAvailabilityChecker.Check(ctx, book); err != nil {
		return NewBorrowBookResponseFailure(err.Error()), nil
	}

	// STEP 6: Create loan entity
	loanID := entity.NewLoanID()
	borrowedAt := time.Now()
	dueDate := borrowedAt.Add(14 * 24 * time.Hour)

	loan := entity.NewLoan(loanID, user.GetID(), book.GetID(), borrowedAt, dueDate)

	// STEP 7: Save loan (THE ONLY WRITE - Single Source of Truth!)
	if err := uc.loanRepository.Save(ctx, loan); err != nil {
		return NewBorrowBookResponseFailure("failed to save loan"), nil
	}

	// ❌ REMOVED: book.MarkAsBorrowed() - No longer needed!
	// ❌ REMOVED: bookRepository.Save(book) - No dual source!

	// STEP 8: Return success response
	return NewBorrowBookResponseSuccess(loan.GetID().String(), dueDate), nil
}