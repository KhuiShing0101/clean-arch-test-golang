package repository

import (
	"testing"
	"library-management/internal/domain/entity"
	"library-management/internal/domain/valueobject"
)

// TestInMemoryUserRepository_Save tests saving users
// InMemoryUserRepository の保存機能をテスト
func TestInMemoryUserRepository_Save(t *testing.T) {
	// Create repository instance
	// リポジトリインスタンスを作成
	repo := NewInMemoryUserRepository()

	// Create test user (ID auto-generated)
	// テストユーザーを作成（ID自動生成）
	user := entity.NewUser("John Doe", "john@example.com")

	// Save user
	// ユーザーを保存
	err := repo.Save(user)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Retrieve and verify
	// 取得して検証
	found, err := repo.FindById(user.Id())
	if err != nil {
		t.Errorf("FindById failed: %v", err)
	}
	if found == nil {
		t.Fatal("User not found after save")
	}
	if found.Email() != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", found.Email())
	}
}

// TestInMemoryUserRepository_FindById tests finding users by ID
// IDによるユーザー検索をテスト
func TestInMemoryUserRepository_FindById(t *testing.T) {
	repo := NewInMemoryUserRepository()

	t.Run("ReturnsNilWhenNotFound", func(t *testing.T) {
		// Try to find non-existent user
		// 存在しないユーザーを検索
		randomId, _ := valueobject.NewUserId("99999999")
		found, err := repo.FindById(randomId)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if found != nil {
			t.Error("Expected nil, got user")
		}
	})

	t.Run("ReturnsUserWhenExists", func(t *testing.T) {
		// Create and save user
		// ユーザーを作成して保存
		user := entity.NewUser("Test User", "test@example.com")
		repo.Save(user)

		// Find by ID
		// IDで検索
		found, err := repo.FindById(user.Id())

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if found == nil {
			t.Fatal("User should be found")
		}
		if found.Name() != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", found.Name())
		}
	})
}

// TestInMemoryUserRepository_FindByEmail tests finding users by email
// メールアドレスによるユーザー検索をテスト
func TestInMemoryUserRepository_FindByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	t.Run("ReturnsNilWhenNotFound", func(t *testing.T) {
		// Search for non-existent email
		// 存在しないメールアドレスを検索
		found, err := repo.FindByEmail("nonexistent@example.com")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if found != nil {
			t.Error("Expected nil, got user")
		}
	})

	t.Run("ReturnsUserWhenExists", func(t *testing.T) {
		// Create and save user
		// ユーザーを作成して保存
		user := entity.NewUser("Email User", "email@example.com")
		repo.Save(user)

		// Find by email
		// メールで検索
		found, err := repo.FindByEmail("email@example.com")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if found == nil {
			t.Fatal("User should be found")
		}
		if found.Name() != "Email User" {
			t.Errorf("Expected name 'Email User', got '%s'", found.Name())
		}
	})
}

// TestInMemoryUserRepository_Update tests updating existing users
// 既存ユーザーの更新をテスト
func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Create and save original user
	// 元のユーザーを作成して保存
	user := entity.NewUser("Original", "original@example.com")
	repo.Save(user)

	// Modify and save again (update)
	// 変更して再保存（更新）
	suspended := user.Suspend()
	repo.Save(suspended)

	// Verify update
	// 更新を確認
	found, _ := repo.FindById(user.Id())
	if found == nil {
		t.Fatal("User should exist after update")
	}
	if found.Status() != entity.UserStatusSuspended {
		t.Errorf("Expected status SUSPENDED, got %s", found.Status())
	}
}