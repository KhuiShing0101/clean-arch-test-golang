package usecase

import (
	"testing"
	"library-management/internal/domain/entity"
	"library-management/internal/domain/repository"
	"library-management/internal/domain/valueobject"
)

// Mock repository for testing - Only 3 methods (simplified)
// テスト用モックリポジトリ - 3メソッドのみ（簡略化）
type mockUserRepository struct {
	users map[string]*entity.User
}

// newMockUserRepository creates a new mock repository
// 新しいモックリポジトリを作成
func newMockUserRepository() repository.UserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

// Save stores a user in memory
// ユーザーをメモリに保存
func (m *mockUserRepository) Save(u *entity.User) error {
	m.users[u.Id().Value()] = u
	return nil
}

// FindById retrieves a user by ID
// IDでユーザーを取得
func (m *mockUserRepository) FindById(id *valueobject.UserId) (*entity.User, error) {
	u, exists := m.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return u, nil
}

// FindByEmail retrieves a user by email (for duplicate check)
// メールでユーザーを取得（重複チェック用）
func (m *mockUserRepository) FindByEmail(email string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, nil
}

// TestCreateUserUseCase tests the CreateUserUseCase
// CreateUserUseCaseをテスト
func TestCreateUserUseCase(t *testing.T) {
	// Test 1: Successfully create a new user
	// テスト1: 新規ユーザーの作成に成功
	t.Run("Success", func(t *testing.T) {
		// Setup: Create mock repository and use case
		// セットアップ: モックリポジトリとユースケースを作成
		repo := newMockUserRepository()
		uc := NewCreateUserUseCase(repo)

		// Execute: Create new user
		// 実行: 新規ユーザーを作成
		input := CreateUserInput{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		output, err := uc.Execute(input)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Assert: Verify all output fields
		// 検証: すべての出力フィールドを確認
		if output.Name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got '%s'", output.Name)
		}
		if output.Email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got '%s'", output.Email)
		}
		if output.Status != "active" {
			t.Errorf("Expected status 'active', got '%s'", output.Status)
		}
		if output.CurrentBorrowCount != 0 {
			t.Errorf("Expected borrow count 0, got %d", output.CurrentBorrowCount)
		}
		if output.OverdueFees != 0 {
			t.Errorf("Expected overdue fees 0, got %.2f", output.OverdueFees)
		}
	})

	// Test 2: Reject duplicate email
	// テスト2: 重複メールを拒否
	t.Run("DuplicateEmail", func(t *testing.T) {
		repo := newMockUserRepository()
		uc := NewCreateUserUseCase(repo)

		// First: Create initial user
		// まず: 最初のユーザーを作成
		input1 := CreateUserInput{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		_, err := uc.Execute(input1)
		if err != nil {
			t.Fatalf("Failed to create first user: %v", err)
		}

		// Second: Try to create user with same email
		// 次に: 同じメールでユーザー作成を試行
		input2 := CreateUserInput{
			Name:  "Jane Doe",
			Email: "john@example.com", // Duplicate / 重複
		}
		_, err = uc.Execute(input2)
		// Should get error for duplicate email
		// 重複メールでエラーを取得するべき
		if err == nil {
			t.Error("Expected error for duplicate email")
		}
	})
}