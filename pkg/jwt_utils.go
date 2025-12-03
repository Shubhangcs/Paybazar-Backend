package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtils struct{}

func (JwtUtils) GenerateToken(data interface{}, expiry time.Duration) (string, error) {
	// JWT secret key
	var key []byte = []byte("agromart@2025")

	// Creating a jwt
	var claims = jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(expiry).Unix(),
	}
	tokenSign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Obtaining the signed string
	token, err := tokenSign.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	// Returning the token
	return token, nil
}

func (JwtUtils) GenerateTokenForExternalAPI(reqid string) (string, error) {
	// JWT secret key
	var key []byte = []byte("UTA5U1VEQXdNREF4VFZSSmVrNUVWVEpPZWxVd1RuYzlQUT09")

	// Creating a jwt
	var claims = jwt.MapClaims{
		"partnerId": "CORP00001",
		"reqid":     reqid,
		"timestamp": time.Now().Unix(),
	}
	tokenSign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Obtaining the signed string
	token, err := tokenSign.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	// Returning the token
	return token, nil
}
