package entity

import (
	"errors"
	"time"
	"library-management/internal/domain/valueobject"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	id                 *valueobject.UserId
	name               string
	email              string
	status             UserStatus
	currentBorrowCount int
	overdueFees        float64
	createdAt          time.Time
}

const MaxBorrowLimit = 5

// NewUser creates a new active user
func NewUser(name, email string) *User {
	return &User{
		id:                 valueobject.GenerateUserId(),
		name:               name,
		email:              email,
		status:             UserStatusActive,
		currentBorrowCount: 0,
		overdueFees:        0,
		createdAt:          time.Now(),
	}
}

// ReconstructUser rebuilds user from persistence
func ReconstructUser(
	id *valueobject.UserId,
	name, email string,
	status UserStatus,
	currentBorrowCount int,
	overdueFees float64,
	createdAt time.Time,
) *User {
	return &User{id, name, email, status, currentBorrowCount, overdueFees, createdAt}
}

func (u *User) CanBorrow() bool {
	if u.status == UserStatusSuspended {
		return false
	}
	if u.currentBorrowCount >= MaxBorrowLimit {
		return false
	}
	if u.overdueFees > 0 {
		return false
	}
	return true
}

func (u *User) BorrowBook() (*User, error) {
	if !u.CanBorrow() {
		return nil, errors.New("user cannot borrow more books")
	}
	return &User{
		id:                 u.id,
		name:               u.name,
		email:              u.email,
		status:             u.status,
		currentBorrowCount: u.currentBorrowCount + 1,
		overdueFees:        u.overdueFees,
		createdAt:          u.createdAt,
	}, nil
}

func (u *User) ReturnBook() (*User, error) {
	if u.currentBorrowCount == 0 {
		return nil, errors.New("user has no books to return")
	}
	return &User{
		id:                 u.id,
		name:               u.name,
		email:              u.email,
		status:             u.status,
		currentBorrowCount: u.currentBorrowCount - 1,
		overdueFees:        u.overdueFees,
		createdAt:          u.createdAt,
	}, nil
}

func (u *User) AddOverdueFee(amount float64) (*User, error) {
	if amount < 0 {
		return nil, errors.New("overdue fee cannot be negative")
	}
	return &User{
		id:                 u.id,
		name:               u.name,
		email:              u.email,
		status:             u.status,
		currentBorrowCount: u.currentBorrowCount,
		overdueFees:        u.overdueFees + amount,
		createdAt:          u.createdAt,
	}, nil
}

func (u *User) Suspend() *User {
	return &User{
		id:                 u.id,
		name:               u.name,
		email:              u.email,
		status:             UserStatusSuspended,
		currentBorrowCount: u.currentBorrowCount,
		overdueFees:        u.overdueFees,
		createdAt:          u.createdAt,
	}
}

// Getters
func (u *User) Id() *valueobject.UserId { return u.id }
func (u *User) Name() string { return u.name }
func (u *User) Email() string { return u.email }
func (u *User) Status() UserStatus { return u.status }
func (u *User) CurrentBorrowCount() int { return u.currentBorrowCount }
func (u *User) OverdueFees() float64 { return u.overdueFees }