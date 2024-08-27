package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ananikitina/notes-rest/internal/domain"
	"github.com/ananikitina/notes-rest/internal/models"
	"github.com/ananikitina/notes-rest/internal/usecases"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
	jwtService  domain.JWTServiceInterface
}

func NewUserHandler(userUseCase *usecases.UserUseCase, jwtService domain.JWTServiceInterface) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, jwtService: jwtService}
}

// RegisterHandler creates new user.
func (u *UserHandler) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		//default role "user" if empty
		if user.Role == "" {
			user.Role = "user"
		}

		if err := u.userUseCase.Register(&user); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	}
}

// LoginHandler authenticates the user.
func (u *UserHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cred struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		user, err := u.userUseCase.Authenticate(cred.Email, cred.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		token, err := u.jwtService.GenerateToken(user.ID, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
