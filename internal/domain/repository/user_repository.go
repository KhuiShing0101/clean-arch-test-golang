package repository

import "library-management/internal/domain/entity"
import "library-management/internal/domain/valueobject"

type UserRepository interface {
	Save(user *entity.User) error
	FindById(id *valueobject.UserId) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}