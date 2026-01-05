package usecase

import (
	"errors"
	"library-management/internal/domain/entity"
	"library-management/internal/domain/repository"
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
	userRepository repository.UserRepository
}

func NewCreateUserUseCase(repo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	// Check duplicate email
	existingUser, _ := uc.userRepository.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user
	user := entity.NewUser(input.Name, input.Email)

	// Save to repository
	if err := uc.userRepository.Save(user); err != nil {
		return nil, err
	}

	// Return DTO
	return &CreateUserOutput{
		Id:                 user.Id().Value(),
		Name:               user.Name(),
		Email:              user.Email(),
		Status:             string(user.Status()),
		CurrentBorrowCount: user.CurrentBorrowCount(),
		OverdueFees:        user.OverdueFees(),
	}, nil
}