package valueobject

import (
	"testing"
)

// TestUserId_Create tests creating a valid UserId
// 有効なUserIdの作成をテスト
func TestUserId_Create(t *testing.T) {
	// Valid 8-digit ID
	// 有効な8桁のID
	userId, err := NewUserId("12345678")
	if err != nil {
		t.Fatalf("Failed to create valid UserId: %v", err)
	}
	if userId.Value() != "12345678" {
		t.Errorf("Expected ID '12345678', got '%s'", userId.Value())
	}
}

// TestUserId_ValidateFormat tests validation of UserId format
// UserIdフォーマットの検証をテスト
func TestUserId_ValidateFormat(t *testing.T) {
	// Invalid: not 8 digits
	// 無効: 8桁ではない
	_, err := NewUserId("123")
	if err == nil {
		t.Error("Expected error for non-8-digit UserId")
	}

	// Invalid: contains letters
	// 無効: 文字を含む
	_, err = NewUserId("1234abcd")
	if err == nil {
		t.Error("Expected error for non-numeric UserId")
	}
}

// TestUserId_Generate tests random UserId generation
// ランダムなUserId生成をテスト
func TestUserId_Generate(t *testing.T) {
	userId := GenerateUserId()
	if userId == nil {
		t.Fatal("GenerateUserId returned nil")
	}
	if len(userId.Value()) != 8 {
		t.Errorf("Generated UserId should be 8 digits, got %d", len(userId.Value()))
	}
}

// TestUserId_Equals tests value equality comparison
// 値の等価性比較をテスト
func TestUserId_Equals(t *testing.T) {
	userId1, _ := NewUserId("12345678")
	userId2, _ := NewUserId("12345678")
	userId3, _ := NewUserId("87654321")

	// Same values should be equal
	// 同じ値は等しいべき
	if !userId1.Equals(userId2) {
		t.Error("UserId with same value should be equal")
	}

	// Different values should not be equal
	// 異なる値は等しくないべき
	if userId1.Equals(userId3) {
		t.Error("UserId with different values should not be equal")
	}
}