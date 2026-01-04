package valueobject

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type UserId struct {
	value string
}

// NewUserId creates a UserId with validation (8 digits)
func NewUserId(value string) (*UserId, error) {
	matched, _ := regexp.MatchString(`^\d{8}$`, value)
	if !matched {
		return nil, errors.New("UserId must be exactly 8 digits. Got: " + value)
	}
	return &UserId{value: value}, nil
}

// GenerateUserId creates a random 8-digit UserId
func GenerateUserId() *UserId {
	rand.Seed(time.Now().UnixNano())
	randomId := rand.Intn(90000000) + 10000000
	return &UserId{value: strconv.Itoa(randomId)}
}

func (u *UserId) Value() string {
	return u.value
}

func (u *UserId) Equals(other *UserId) bool {
	if other == nil {
		return false
	}
	return u.value == other.value
}

func (u *UserId) String() string {
	return u.value
}