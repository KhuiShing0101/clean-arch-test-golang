package repository

import (
	"library-management/internal/domain/user"
)

type InMemoryUserRepository struct {
	users map[string]*user.User
}

func NewInMemoryUserRepository() user.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*user.User),
	}
}

func (r *InMemoryUserRepository) Save(u *user.User) error {
	r.users[u.Id().Value()] = u
	return nil
}

func (r *InMemoryUserRepository) FindById(id *user.UserId) (*user.User, error) {
	u, exists := r.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Delete(id *user.UserId) error {
	delete(r.users, id.Value())
	return nil
}

func (r *InMemoryUserRepository) FindUsersWithOverdueFees() ([]*user.User, error) {
	var users []*user.User
	for _, u := range r.users {
		if u.OverdueFees() > 0 {
			users = append(users, u)
		}
	}
	return users, nil
}