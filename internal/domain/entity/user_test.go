package entity

import (
	"testing"
	"time"
	"library-management/internal/domain/valueobject"
)

// TestUser tests all User entity business logic
// Userエンティティの全ビジネスロジックをテスト
func TestUser(t *testing.T) {
	// Test 1: New user can borrow books
	// テスト1: 新規ユーザーは本を借りられる
	t.Run("CanBorrow", func(t *testing.T) {
		u := NewUser("John Doe", "john@example.com")

		// Check if new user can borrow
		// 新規ユーザーが借りられるかチェック
		if !u.CanBorrow() {
			t.Error("New user should be able to borrow")
		}

		// Verify user properties
		// ユーザープロパティを検証
		if u.Name() != "John Doe" {
			t.Errorf("Expected name 'John Doe', got '%s'", u.Name())
		}
		if u.Email() != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got '%s'", u.Email())
		}
		if u.Status() != UserStatusActive {
			t.Errorf("Expected status 'active', got '%s'", u.Status())
		}
		if u.CurrentBorrowCount() != 0 {
			t.Errorf("Expected borrow count 0, got %d", u.CurrentBorrowCount())
		}
	})

	// Test 2: Borrowing a book increments count
	// テスト2: 本を借りるとカウントが増える
	t.Run("BorrowBook", func(t *testing.T) {
		u := NewUser("John Doe", "john@example.com")

		u2, err := u.BorrowBook()
		if err != nil {
			t.Fatalf("Failed to borrow book: %v", err)
		}

		// New instance has incremented count
		// 新しいインスタンスはカウントが増加
		if u2.CurrentBorrowCount() != 1 {
			t.Errorf("Expected borrow count 1, got %d", u2.CurrentBorrowCount())
		}

		// Original user unchanged (immutability)
		// 元のユーザーは不変（不変性）
		if u.CurrentBorrowCount() != 0 {
			t.Error("Original user should remain unchanged")
		}
	})

	// Test 3: Suspended users cannot borrow
	// テスト3: 停止中のユーザーは借りられない
	t.Run("CannotBorrowWhenSuspended", func(t *testing.T) {
		userId, _ := valueobject.NewUserId("12345678")

		// Create suspended user
		// 停止中のユーザーを作成
		u := ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			UserStatusSuspended,
			0,
			0,
			time.Now(),
		)

		if u.CanBorrow() {
			t.Error("Suspended user should not be able to borrow")
		}
	})

	// Test 4: Cannot borrow at max limit (5 books)
	// テスト4: 最大制限（5冊）では借りられない
	t.Run("CannotBorrowAtMaxLimit", func(t *testing.T) {
		userId, _ := valueobject.NewUserId("12345678")

		// User at max borrow limit
		// 最大貸出制限のユーザー
		u := ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			UserStatusActive,
			MaxBorrowLimit, // 5 books / 5冊
			0,
			time.Now(),
		)

		if u.CanBorrow() {
			t.Error("User should not be able to borrow (max loans reached)")
		}
	})

	// Test 5: Cannot borrow with overdue fees
	// テスト5: 延滞料金があると借りられない
	t.Run("CannotBorrowWithOverdueFees", func(t *testing.T) {
		userId, _ := valueobject.NewUserId("12345678")

		// User with overdue fees
		// 延滞料金があるユーザー
		u := ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			UserStatusActive,
			0,
			10.50, // ¥10.50 in fees / 延滞料金¥10.50
			time.Now(),
		)

		if u.CanBorrow() {
			t.Error("User with overdue fees should not be able to borrow")
		}
	})

	// Test 6: Fee management (add/pay)
	// テスト6: 料金管理（追加/支払い）
	t.Run("FeeManagement", func(t *testing.T) {
		u := NewUser("John Doe", "john@example.com")

		// Add overdue fee
		// 延滞料金を追加
		u2, err := u.AddOverdueFee(5.00)
		if err != nil {
			t.Fatalf("Failed to add overdue fee: %v", err)
		}

		if u2.OverdueFees() != 5.00 {
			t.Errorf("Expected overdue fees 5.00, got %.2f", u2.OverdueFees())
		}

		// Negative fee should error
		// 負の料金はエラーになるべき
		_, err = u.AddOverdueFee(-1.00)
		if err == nil {
			t.Error("Expected error for negative fee")
		}
	})

	// Test 7: Immutability pattern (state changes return new instances)
	// テスト7: 不変性パターン（状態変更は新しいインスタンスを返す）
	t.Run("Immutability", func(t *testing.T) {
		userId, _ := valueobject.NewUserId("12345678")
		// User with 2 borrowed books
		// 2冊借りているユーザー
		u := ReconstructUser(
			userId,
			"John Doe",
			"john@example.com",
			UserStatusActive,
			2,
			0,
			time.Now(),
		)

		// Return book - should get new instance
		// 本を返却 - 新しいインスタンスを取得するべき
		u2, err := u.ReturnBook()
		if err != nil {
			t.Fatalf("Failed to return book: %v", err)
		}

		// New instance has decremented count
		// 新しいインスタンスはカウント減少
		if u2.CurrentBorrowCount() != 1 {
			t.Errorf("Expected borrow count 1, got %d", u2.CurrentBorrowCount())
		}

		// Original instance unchanged
		// 元のインスタンスは変更なし
		if u.CurrentBorrowCount() != 2 {
			t.Error("Original user should remain unchanged")
		}
	})
}