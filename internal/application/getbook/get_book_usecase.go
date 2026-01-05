package getbook

import "library-management/internal/application/query"

// GetBookUseCase - GetBook Use Case (CQRS Query Side)
//
// This use case is very thin - it just delegates to QueryService.
// All the complexity is in BookQueryService!
type GetBookUseCase struct {
	bookQueryService query.BookQueryService
}

func NewGetBookUseCase(bookQueryService query.BookQueryService) *GetBookUseCase {
	return &GetBookUseCase{
		bookQueryService: bookQueryService,
	}
}

func (uc *GetBookUseCase) Execute(request *GetBookRequest) *GetBookResponse {
	// Simply delegate to QueryService!
	bookReadModel, err := uc.bookQueryService.GetBookById(request.BookId)

	if err != nil {
		return NewFailureResponse("An unexpected error occurred")
	}

	if bookReadModel == nil {
		return NewNotFoundResponse()
	}

	return NewSuccessResponse(bookReadModel)
}