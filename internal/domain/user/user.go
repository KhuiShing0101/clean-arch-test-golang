package user

import "errors"

// Value Objects
type UserId string
type UserName string
type Email string

const MaxLoans = 5

type User struct {
	id               UserId
	name             UserName
	email            Email
	currentLoanCount int
}

// CanBorrow - validates loan eligibility
func (u *User) CanBorrow() bool {
	return u.currentLoanCount < MaxLoans
}

// RecordLoan - records new loan
func (u *User) RecordLoan() error {
	if !u.CanBorrow() {
		return errors.New("user has reached loan limit")
	}
	u.currentLoanCount++
	return nil
}

// RecordReturn - records return
func (u *User) RecordReturn() error {
	if u.currentLoanCount == 0 {
		return errors.New("no active loans to return")
	}
	u.currentLoanCount--
	return nil
}

// Getters
func (u *User) GetId() UserId {
	return u.id
}

func (u *User) GetName() UserName {
	return u.name
}

func (u *User) GetEmail() Email {
	return u.email
}

func (u *User) GetCurrentLoanCount() int {
	return u.currentLoanCount
}

// Constructor
func NewUser(id UserId, name UserName, email Email) *User {
	return &User{
		id:               id,
		name:             name,
		email:            email,
		currentLoanCount: 0,
	}
}

// Repository Interface (Domain Layer)
type IUserRepository interface {
	FindById(id UserId) (*User, error)
	Save(user *User) error
}