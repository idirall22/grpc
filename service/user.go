package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Username     string
	HashPassword string
	Role         string
}

// NewUser create a new user
func NewUser(username, password, role string) (*User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Coulod not hash password %v", err)
	}
	return &User{
		Username:     username,
		HashPassword: string(hashPassword),
		Role:         role,
	}, nil
}

// IsCorrectPassword check if password is correct
func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(password))
	return err == nil
}

// Clone clone user
func (u *User) Clone() *User {
	return &User{
		Username:     u.Username,
		HashPassword: u.HashPassword,
		Role:         u.Role,
	}
}
