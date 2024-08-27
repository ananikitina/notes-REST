package services

import (
	"errors"
	"time"

	"github.com/ananikitina/notes-rest/internal/config"
	"github.com/ananikitina/notes-rest/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

// JWTService struct for generating and validating JWT tokens.
type JWTService struct {
	secret string
	issuer string
}

// NewJWTService creates a new JWT service with given configuration.
func NewJWTService(cfg *config.Config) domain.JWTServiceInterface {
	return &JWTService{secret: cfg.JWTSecret, issuer: "notes-rest"}
}

// GenerateToken creates a JWT token for a user.
func (j *JWTService) GenerateToken(userID int, userRole string) (string, error) {
	claims := &domain.Claims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken validates the JWT token and returns the claims.
func (j *JWTService) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
