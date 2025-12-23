package book

import "errors"

// Book Entity - represents a book in our library
type Book struct {
  id              *BookId
  title           string
  author          string
  isbn            *ISBN
  totalCopies     int
  availableCopies int
}

func NewBook(
  id *BookId,
  title string,
  author string,
  isbn *ISBN,
  totalCopies int,
) (*Book, error) {
  // Validation in constructor (Fail Fast)
  if len(title) == 0 {
    return nil, errors.New("book title cannot be empty")
  }
  if len(author) == 0 {
    return nil, errors.New("book author cannot be empty")
  }
  if totalCopies < 1 {
    return nil, errors.New("total copies must be at least 1")
  }

  return &Book{
    id:              id,
    title:           title,
    author:          author,
    isbn:            isbn,
    totalCopies:     totalCopies,
    availableCopies: totalCopies, // Initially all copies available
  }, nil
}

// Getters
func (b *Book) GetId() *BookId {
  return b.id
}

func (b *Book) GetTitle() string {
  return b.title
}

func (b *Book) GetAuthor() string {
  return b.author
}

func (b *Book) GetISBN() *ISBN {
  return b.isbn
}

func (b *Book) GetTotalCopies() int {
  return b.totalCopies
}

func (b *Book) GetAvailableCopies() int {
  return b.availableCopies
}

// Business Logic Methods
func (b *Book) IsAvailable() bool {
  return b.availableCopies > 0
}

func (b *Book) BorrowCopy() error {
  if !b.IsAvailable() {
    return errors.New("no copies available for book: " + b.title)
  }
  b.availableCopies--
  return nil
}

func (b *Book) ReturnCopy() error {
  if b.availableCopies >= b.totalCopies {
    return errors.New("cannot return more copies than total")
  }
  b.availableCopies++
  return nil
}

// Mutable state change methods
func (b *Book) UpdateTitle(newTitle string) error {
  if len(newTitle) == 0 {
    return errors.New("book title cannot be empty")
  }
  b.title = newTitle
  return nil
}

func (b *Book) UpdateAuthor(newAuthor string) error {
  if len(newAuthor) == 0 {
    return errors.New("book author cannot be empty")
  }
  b.author = newAuthor
  return nil
}