package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenDuration is how long an access token is valid.
const AccessTokenDuration = 15 * time.Minute

// Manager creates and validates JWT access tokens.
type Manager struct {
	secret []byte
}

// New creates a new Manager with the given signing secret.
func New(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

// Claims are the JWT claims for an access token.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed access token for the given user ID.
func (manager *Manager) GenerateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(manager.secret)
}

// ParseAccessToken validates the token and returns the user ID.
func (manager *Manager) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return manager.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}
