package usecases

import (
	"errors"

	"github.com/ananikitina/notes-rest/internal/models"
	"github.com/ananikitina/notes-rest/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase represents the business logic for users.
type UserUseCase struct {
	userRepo *repository.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase.
func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) Register(user *models.User) error {
	// Check if user already exists
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}
	return u.userRepo.CreateUser(user)
}

func (u *UserUseCase) Authenticate(email, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}
