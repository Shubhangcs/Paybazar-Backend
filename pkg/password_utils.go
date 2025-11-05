package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordUtils struct{}

func (*PasswordUtils) HashPassword(password string) (string, error) {
	// Encrypt password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	// Returning the encrypted password in string format
	return string(hashPass), nil
}

func (*PasswordUtils) VerifyPassword(hashedPassword string, normalPassword string) error {
	// Compare the encrypted password with normal password and return the result
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(normalPassword))
}