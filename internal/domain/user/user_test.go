package user_test

import (
	"testing"
	"library-management/internal/domain/user"
)

func TestCanCreateValidUser(t *testing.T) {
	// Create Value Object
	userId, _ := user.NewUserId("user-123")

	u, err := user.NewUser(userId, "John Doe", "john@example.com", 0, false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if u.GetId().GetValue() != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", u.GetId().GetValue())
	}
	if u.GetName() != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", u.GetName())
	}
	if u.GetEmail() != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", u.GetEmail())
	}
}

func TestInvalidEmailReturnsError(t *testing.T) {
	userId, _ := user.NewUserId("user-123")

	_, err := user.NewUser(userId, "John Doe", "invalid-email", 0, false)
	if err == nil {
		t.Error("Expected error for invalid email")
	}
}

func TestCannotBorrowMoreWhenMaxLoansReached(t *testing.T) {
	userId, _ := user.NewUserId("user-123")
	u, _ := user.NewUser(userId, "John Doe", "john@example.com", 5, false)

	if u.CanBorrowMore() {
		t.Error("User should not be able to borrow more (max loans reached)")
	}
}

func TestCanBorrowMoreWhenUnderLimit(t *testing.T) {
	userId, _ := user.NewUserId("user-123")
	u, _ := user.NewUser(userId, "John Doe", "john@example.com", 2, false)

	if !u.CanBorrowMore() {
		t.Error("User should be able to borrow more")
	}
}
