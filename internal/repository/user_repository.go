package repository

import (
	"database/sql"

	"github.com/ananikitina/notes-rest/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// userRepository handles user-related database operations.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new user repository with the given database connection.
func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// CreateUser adds a new user to the database.
func (u *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = u.DB.Exec(
		"INSERT INTO users (email,password,role) VALUES ($1, $2,$3)",
		user.Email, user.Password, user.Role)
	return err
}

// GetUserByEmail retrieves a user by their email.
func (u *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := u.DB.QueryRow(
		"SELECT id, email, password, role FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
