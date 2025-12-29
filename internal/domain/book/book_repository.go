package book

type BookRepository interface {
	Save(book *Book) error
	FindById(id *BookId) (*Book, error)
	FindByISBN(isbn *ISBN) (*Book, error)
	FindAll() ([]*Book, error)
	Delete(id *BookId) error
}