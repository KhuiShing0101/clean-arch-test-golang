package borrowbook

import (
	"library-management/internal/domain/user"
	"library-management/internal/domain/book"
	"library-management/internal/domain/loan"
	"library-management/internal/domain/shared"
)

type BorrowBookUseCase struct {
	userRepo    user.IUserRepository
	bookRepo    book.IBookRepository
	loanRepo    loan.ILoanRepository
	policyService *loan.BorrowingPolicyService
	txManager   shared.TransactionManager
}

func NewBorrowBookUseCase(
	userRepo user.IUserRepository,
	bookRepo book.IBookRepository,
	loanRepo loan.ILoanRepository,
	policyService *loan.BorrowingPolicyService,
	txManager shared.TransactionManager,
) *BorrowBookUseCase {
	return &BorrowBookUseCase{
		userRepo:      userRepo,
		bookRepo:      bookRepo,
		loanRepo:      loanRepo,
		policyService: policyService,
		txManager:     txManager,
	}
}

// Execute - borrows a book - coordinates multiple entities within transaction
func (uc *BorrowBookUseCase) Execute(req *BorrowBookRequest) (*BorrowBookResponse, error) {
	var response *BorrowBookResponse

	err := uc.txManager.RunInTransaction(func() error {
		// 1. Fetch user entity
		u, err := uc.userRepo.FindById(req.UserId)
		if err != nil {
			return &UserNotFoundError{UserId: req.UserId}
		}

		// 2. Fetch book entity
		b, err := uc.bookRepo.FindById(req.BookId)
		if err != nil {
			return &BookNotFoundError{BookId: req.BookId}
		}

		// 3. Validate business rules (domain service)
		if err := uc.policyService.CanBorrow(u, b); err != nil {
			return err
		}

		// 4. Create new loan entity via repository (DIP - repository creates domain entities)
		l, err := uc.loanRepo.CreateLoan(u.GetId(), b.GetId())
		if err != nil {
			return err
		}

		// 5. Update user and book state
		if err := u.RecordLoan(); err != nil {
			return err
		}
		if err := b.MarkAsBorrowed(); err != nil {
			return err
		}

		// 6. Save all entities (within transaction)
		if err := uc.userRepo.Save(u); err != nil {
			return err
		}
		if err := uc.bookRepo.Save(b); err != nil {
			return err
		}
		if err := uc.loanRepo.Save(l); err != nil {
			return err
		}

		// 7. Create response
		response = NewBorrowBookResponse(l, b.GetTitle().String())
		return nil
	})

	if err != nil {
		return nil, err
	}
	return response, nil
}