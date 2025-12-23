package book

import "errors"

// WHY a struct? Type safety - can't accidentally use string as BookId
// WHY unexported field? Immutability - no external code can change this
type BookId struct {
  value string
}

// WHY constructor function? "Fail fast" - catch errors early
// After this runs successfully, object is GUARANTEED valid
func NewBookId(value string) (*BookId, error) {
  // WHY check empty? Business rule: IDs must exist
  // Better to fail here than get empty ID in database
  if len(value) == 0 {
    return nil, errors.New("BookId cannot be empty")
  }

  // WHY return pointer? Conventional in Go for objects
  // Value is immutable because field is unexported
  return &BookId{value: value}, nil
}

// WHY getter? Need to access value without exposing internals
// Returns primitive for database/API use
func (b *BookId) GetValue() string {
  return b.value
}

// WHY Equals method? Value Objects compare by VALUE not identity
// Two BookId("123") are "equal" even if different pointers
func (b *BookId) Equals(other *BookId) bool {
  return b.value == other.value
}
// WHY String() method? Implements fmt.Stringer interface
// fmt.Println(bookId) works naturally
// String returns the string representation of the BookId.
// This method implements the fmt.Stringer interface, allowing BookId
// to be easily converted to a human-readable string format.
func (b *BookId) String() string {
	return b.value
}
