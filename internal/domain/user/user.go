package user

import (
  "errors"
  "regexp"
)

// User Entity - represents a library user
type User struct {
  id               *UserId
  name             string
  email            string
  currentLoanCount int
  hasOverdueBooks  bool
}

const MaxLoans = 5

func NewUser(
  id *UserId,
  name string,
  email string,
  currentLoanCount int,
  hasOverdueBooks bool,
) (*User, error) {
  if len(name) == 0 {
    return nil, errors.New("name cannot be empty")
  }

  if !isValidEmail(email) {
    return nil, errors.New("invalid email format")
  }

  if currentLoanCount < 0 {
    return nil, errors.New("loan count cannot be negative")
  }

  return &User{
    id:               id,
    name:             name,
    email:            email,
    currentLoanCount: currentLoanCount,
    hasOverdueBooks:  hasOverdueBooks,
  }, nil
}

func isValidEmail(email string) bool {
  re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
  return re.MatchString(email)
}

func (u *User) CanBorrowMore() bool {
  return u.currentLoanCount < MaxLoans
}

func (u *User) IncrementLoanCount() {
  u.currentLoanCount++
}

func (u *User) DecrementLoanCount() {
  if u.currentLoanCount > 0 {
    u.currentLoanCount--
  }
}

func (u *User) MarkAsHavingOverdueBooks() {
  u.hasOverdueBooks = true
}

func (u *User) ClearOverdueBooks() {
  u.hasOverdueBooks = false
}

// Getters
func (u *User) GetId() *UserId {
  return u.id
}

func (u *User) GetName() string {
  return u.name
}

func (u *User) GetEmail() string {
  return u.email
}

func (u *User) GetCurrentLoanCount() int {
  return u.currentLoanCount
}

func (u *User) HasOverdueBooks() bool {
  return u.hasOverdueBooks
}

func (u *User) GetMaxLoans() int {
  return MaxLoans
}