package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTManager struct
type JWTManager struct {
	Secret        string
	TokenDuration time.Duration
}

// UserClaims struct
type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     string
}

// NewJWTManager create JWTManager
func NewJWTManager(secret string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		Secret:        secret,
		TokenDuration: tokenDuration,
	}
}

// Generate generate a token
func (m *JWTManager) Generate(user *User) (string, error) {
	userClaims := &UserClaims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * m.TokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	return token.SignedString([]byte(m.Secret))
}

// Verify verify token
func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(m.Secret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, fmt.Errorf("Invalid claims: %v", err)
	}
	return claims, nil
}
