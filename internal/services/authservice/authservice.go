package authservice

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	signingMethod *jwt.SigningMethodHMAC
	signingKey    interface{}
}

type AuthServiceInterface interface {
	GenerateToken(c Claimer) (string, error)
}

type Claimer interface {
	GetClaims() map[string]interface{}
}

func New(method *jwt.SigningMethodHMAC, key interface{}) AuthServiceInterface {
	return &AuthService{method, key}
}

func (s AuthService) GenerateToken(c Claimer) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range c.GetClaims() {
		claims[k] = v
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)

	return token.SignedString(s.signingKey)
}

