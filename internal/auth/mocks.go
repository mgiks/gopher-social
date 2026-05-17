package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewMockAuthenticator() Authenticator {
	return MockAuthenticator{}
}

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-aud",
	"sub": 42,
	"exp": time.Now().Add(time.Hour).Unix(),
}

type MockAuthenticator struct{}

func (a MockAuthenticator) GenerateToken(jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret))
}

func (a MockAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
