package authservice

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	signingMethod      jwt.SigningMethod
	signingKey         interface{}
	allowSigningMethod AllowSigningMethod
}

type AuthServiceInterface interface {
	GenerateToken(c Claimer, expiredIn time.Duration) (string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type AllowSigningMethod struct {
	ECDSA   bool
	Ed25519 bool
	HMAC    bool
	RSA     bool
	RSAPSS  bool
}

func (a *AllowSigningMethod) Allowed(k string) {
	ps := reflect.ValueOf(a)
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		f := s.FieldByName(k)
		if f.IsValid() {
			if f.CanSet() {
				if f.Kind() == reflect.Bool {
					f.SetBool(true)
				}
			}
		}
	}
}

type Claimer interface {
	GetClaims() map[string]interface{}
}

func New(method jwt.SigningMethod, key interface{}, allowSigningMethod AllowSigningMethod) AuthServiceInterface {
	return &AuthService{method, key, allowSigningMethod}
}

func (s AuthService) GenerateToken(c Claimer, expiredIn time.Duration) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range c.GetClaims() {
		claims[k] = v
	}

	claims["iat"] = time.Now()
	claims["nbf"] = time.Now()
	claims["exp"] = time.Now().Add(expiredIn).Unix()

	token := jwt.NewWithClaims(s.signingMethod, claims)

	return token.SignedString(s.signingKey)
}

func (s AuthService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		var ok bool
		switch token.Method.(type) {
		case *jwt.SigningMethodECDSA:
			ok = s.allowSigningMethod.ECDSA
		case *jwt.SigningMethodHMAC:
			ok = s.allowSigningMethod.HMAC
		case *jwt.SigningMethodRSA:
			ok = s.allowSigningMethod.RSA
		case *jwt.SigningMethodRSAPSS:
			ok = s.allowSigningMethod.RSAPSS
		}

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
