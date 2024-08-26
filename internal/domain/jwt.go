package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTService defines the contract for JWT operations.
type JWTServiceInterface interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Claims struct represents the JWT claims.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
