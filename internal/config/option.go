package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maetad/baroness-api/internal/services/authservice"
)

type Options struct {
	AppName           string
	ListenAddressHTTP string
	DatabaseHost      string
	DatabaseUser      string
	DatabasePass      string
	DatabaseName      string
	DatabasePort      int
	DatabaseSSLMode   string
	DatabaseTimezone  string
	JWTSigningMethod  jwt.SigningMethod
	JWTSigningKey     []byte
	JWTAllowMethod    authservice.AllowSigningMethod
	JWTExpiredIn      time.Duration
	GCSProjectID      string
	GCSBucket         string
	GoogleAccessID    string
	GoogleCredential  string
	GooglePrivateKey  []byte
}
