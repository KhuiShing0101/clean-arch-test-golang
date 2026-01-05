package getbook

import "library-management/internal/application/query"

// BookDTO - Data Transfer Object for book information
// Maps from BookReadModel to application layer DTO
type BookDTO struct {
	BookId      string
	ISBN        string
	Title       string
	Author      string
	IsAvailable bool
}

// FromReadModel converts BookReadModel to BookDTO
func FromReadModel(readModel *query.BookReadModel) BookDTO {
	return BookDTO{
		BookId:      readModel.ID,
		ISBN:        readModel.ISBN,
		Title:       readModel.Title,
		Author:      readModel.Author,
		IsAvailable: readModel.IsAvailable,
	}
}

// FromReadModels converts array of BookReadModels to BookDTOs
func FromReadModels(readModels []*query.BookReadModel) []BookDTO {
	dtos := make([]BookDTO, len(readModels))
	for i, model := range readModels {
		dtos[i] = FromReadModel(model)
	}
	return dtos
}

// IBookQueryService - Interface for query service used by both GetBook and ListBooks
type IBookQueryService interface {
	GetBook(bookId string) (*BookDTO, error)
	ListBooks(limit int, offset int) ([]BookDTO, int, error)
}

// BookQueryServiceAdapter - Adapter to convert query.BookQueryService to IBookQueryService
type BookQueryServiceAdapter struct {
	queryService query.BookQueryService
}

// NewBookQueryServiceAdapter creates new adapter
func NewBookQueryServiceAdapter(queryService query.BookQueryService) IBookQueryService {
	return &BookQueryServiceAdapter{
		queryService: queryService,
	}
}

// GetBook implements IBookQueryService.GetBook
func (a *BookQueryServiceAdapter) GetBook(bookId string) (*BookDTO, error) {
	readModel, err := a.queryService.GetBookById(bookId)
	if err != nil {
		return nil, err
	}
	if readModel == nil {
		return nil, nil
	}
	dto := FromReadModel(readModel)
	return &dto, nil
}

// ListBooks implements IBookQueryService.ListBooks
func (a *BookQueryServiceAdapter) ListBooks(limit int, offset int) ([]BookDTO, int, error) {
	readModels, total, err := a.queryService.ListBooks(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	dtos := FromReadModels(readModels)
	return dtos, total, nil
}
