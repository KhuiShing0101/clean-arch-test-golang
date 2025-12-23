package loan

import "errors"

type LoanId struct {
  value string
}

func NewLoanId(value string) (*LoanId, error) {
  if len(value) == 0 {
    return nil, errors.New("LoanId cannot be empty")
  }
  return &LoanId{value: value}, nil
}

func (l *LoanId) GetValue() string {
  return l.value
}

func (l *LoanId) Equals(other *LoanId) bool {
  return l.value == other.value
}

func (l *LoanId) String() string {
  return l.value
}