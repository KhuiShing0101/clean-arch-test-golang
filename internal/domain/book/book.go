package book

import "errors"

// Value Objects
type BookId string
type ISBN string
type Title string
type Author string

type BookStatus string

const (
	StatusAvailable BookStatus = "AVAILABLE"
	StatusBorrowed  BookStatus = "BORROWED"
)

type Book struct {
	id     BookId
	isbn   ISBN
	title  Title
	author Author
	status BookStatus
}

// IsAvailable - checks if available for borrowing
func (b *Book) IsAvailable() bool {
	return b.status == StatusAvailable
}

// MarkAsBorrowed - marks as borrowed
func (b *Book) MarkAsBorrowed() error {
	if !b.IsAvailable() {
		return errors.New("book is already borrowed")
	}
	b.status = StatusBorrowed
	return nil
}

// MarkAsAvailable - marks as available
func (b *Book) MarkAsAvailable() error {
	if b.IsAvailable() {
		return errors.New("book is already available")
	}
	b.status = StatusAvailable
	return nil
}

// Getters
func (b *Book) GetId() BookId {
	return b.id
}

func (b *Book) GetISBN() ISBN {
	return b.isbn
}

func (b *Book) GetTitle() Title {
	return b.title
}

func (b *Book) GetAuthor() Author {
	return b.author
}

func (b *Book) GetStatus() BookStatus {
	return b.status
}

// Constructor
func NewBook(id BookId, isbn ISBN, title Title, author Author) *Book {
	return &Book{
		id:     id,
		isbn:   isbn,
		title:  title,
		author: author,
		status: StatusAvailable,
	}
}

// Value object helpers
func (t Title) String() string {
	return string(t)
}

func (a Author) String() string {
	return string(a)
}

// Repository Interface (Domain Layer)
type IBookRepository interface {
	FindById(id BookId) (*Book, error)
	Save(book *Book) error
}