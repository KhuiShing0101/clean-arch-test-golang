package user_test

import (
	"testing"
	"library-management/internal/domain/user"
)

func TestUserId(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		// Valid 8-digit ID
		userId, err := user.NewUserId("12345678")
		if err != nil {
			t.Fatalf("Failed to create valid UserId: %v", err)
		}
		if userId.Value() != "12345678" {
			t.Errorf("Expected ID '12345678', got '%s'", userId.Value())
		}
	})

	t.Run("ValidateFormat", func(t *testing.T) {
		// Invalid: not 8 digits
		_, err := user.NewUserId("123")
		if err == nil {
			t.Error("Expected error for non-8-digit UserId")
		}

		// Invalid: contains letters
		_, err = user.NewUserId("1234abcd")
		if err == nil {
			t.Error("Expected error for non-numeric UserId")
		}
	})

	t.Run("Generate", func(t *testing.T) {
		userId := user.GenerateUserId()
		if userId == nil {
			t.Fatal("GenerateUserId returned nil")
		}
		if len(userId.Value()) != 8 {
			t.Errorf("Generated UserId should be 8 digits, got %d", len(userId.Value()))
		}
	})
}