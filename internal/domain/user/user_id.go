package user

import "errors"

type UserId struct {
  value string
}

func NewUserId(value string) (*UserId, error) {
  if len(value) == 0 {
    return nil, errors.New("UserId cannot be empty")
  }
  return &UserId{value: value}, nil
}

func (u *UserId) GetValue() string {
  return u.value
}

func (u *UserId) Equals(other *UserId) bool {
  return u.value == other.value
}

func (u *UserId) String() string {
  return u.value
}