package authservice

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
}

type AuthServiceInterface interface {
	GenerateToken(c Claimer) (string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type Claimer interface {
	GetClaims() map[string]interface{}
}

func New(method jwt.SigningMethod, key interface{}) AuthServiceInterface {
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

func (s AuthService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
