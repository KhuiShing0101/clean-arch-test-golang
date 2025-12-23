package usecase

import (
	"errors"
	"library-management/internal/domain/user"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	Id                 string
	Name               string
	Email              string
	Status             string
	CurrentBorrowCount int
	OverdueFees        float64
}

type CreateUserUseCase struct {
	userRepository user.UserRepository
}

func NewCreateUserUseCase(repo user.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	// Check duplicate email
	existingUser, _ := uc.userRepository.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user
	u := user.NewUser(input.Name, input.Email)

	// Save to repository
	if err := uc.userRepository.Save(u); err != nil {
		return nil, err
	}

	// Return DTO
	return &CreateUserOutput{
		Id:                 u.Id().Value(),
		Name:               u.Name(),
		Email:              u.Email(),
		Status:             string(u.Status()),
		CurrentBorrowCount: u.CurrentBorrowCount(),
		OverdueFees:        u.OverdueFees(),
	}, nil
}