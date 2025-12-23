package user

type UserRepository interface {
	Save(user *User) error
	FindById(id *UserId) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Delete(id *UserId) error
	FindUsersWithOverdueFees() ([]*User, error)
}