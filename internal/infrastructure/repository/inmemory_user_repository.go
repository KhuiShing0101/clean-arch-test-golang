package repository

import (
	"library-management/internal/domain/entity"
	"library-management/internal/domain/repository"
	"library-management/internal/domain/valueobject"
)

type InMemoryUserRepository struct {
	users map[string]*entity.User
}

func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *InMemoryUserRepository) Save(user *entity.User) error {
	r.users[user.Id().Value()] = user
	return nil
}

func (r *InMemoryUserRepository) FindById(id *valueobject.UserId) (*entity.User, error) {
	user, exists := r.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Email() == email {
			return user, nil
		}
	}
	return nil, nil
}