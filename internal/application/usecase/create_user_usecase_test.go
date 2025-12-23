package usecase_test

import (
	"testing"
	"library-management/internal/application/usecase"
	"library-management/internal/domain/user"
)

// Mock repository for testing
type mockUserRepository struct {
	users map[string]*user.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*user.User),
	}
}

func (m *mockUserRepository) Save(u *user.User) error {
	m.users[u.Id().Value()] = u
	return nil
}

func (m *mockUserRepository) FindById(id *user.UserId) (*user.User, error) {
	u, exists := m.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockUserRepository) Delete(id *user.UserId) error {
	delete(m.users, id.Value())
	return nil
}

func (m *mockUserRepository) FindUsersWithOverdueFees() ([]*user.User, error) {
	var users []*user.User
	for _, u := range m.users {
		if u.OverdueFees() > 0 {
			users = append(users, u)
		}
	}
	return users, nil
}

func TestCreateUserUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := newMockUserRepository()
		uc := usecase.NewCreateUserUseCase(repo)

		input := usecase.CreateUserInput{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		output, err := uc.Execute(input)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

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

		// Verify user was saved to repository
		users, _ := repo.FindAll()
		if len(users) != 1 {
			t.Errorf("Expected 1 user in repository, got %d", len(users))
		}
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		repo := newMockUserRepository()
		uc := usecase.NewCreateUserUseCase(repo)

		// Create first user
		input1 := usecase.CreateUserInput{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		_, err := uc.Execute(input1)
		if err != nil {
			t.Fatalf("Failed to create first user: %v", err)
		}

		// Try to create second user with same email
		input2 := usecase.CreateUserInput{
			Name:  "Jane Doe",
			Email: "john@example.com", // Duplicate email
		}
		_, err = uc.Execute(input2)
		if err == nil {
			t.Error("Expected error for duplicate email")
		}

		// Verify only one user exists
		users, _ := repo.FindAll()
		if len(users) != 1 {
			t.Errorf("Expected 1 user in repository, got %d", len(users))
		}
	})
}
