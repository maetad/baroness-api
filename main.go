package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maetad/baroness-api/internal"
	"github.com/maetad/baroness-api/internal/config"
	"github.com/maetad/baroness-api/internal/services/authservice"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var options config.Options

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := internal.New(ctx, log, options)
	if err != nil {
		log.WithError(err).Fatal("internal.New()")
	}

	go func() {
		// service connections
		if err := svc.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	shutdownOnSignal(svc)

	svc.Close()
}

func init() {
	options = config.Options{
		AppName:           os.Getenv("APP_NAME"),
		ListenAddressHTTP: os.Getenv("LISTEN_ADDRESS_HTTP"),
		DatabaseHost:      os.Getenv("DATABASE_HOST"),
		DatabaseUser:      os.Getenv("DATABASE_USER"),
		DatabasePass:      os.Getenv("DATABASE_PASS"),
		DatabaseName:      os.Getenv("DATABASE_NAME"),
		DatabasePort: func() int {
			i, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
			return i
		}(),
		DatabaseSSLMode: func() string {
			allow := []string{"enable", "require"}
			for _, a := range allow {
				if os.Getenv("DATABASE_SSL_MODE") == a {
					return a
				}
			}

			return "disable"
		}(),
		DatabaseTimezone: os.Getenv("DATABASE_TIMEZONE"),
		JWTSigningMethod: func() jwt.SigningMethod {
			method := os.Getenv("JWT_SIGNING_METHOD")
			switch method {
			case "HS256":
				return jwt.SigningMethodHS256
			case "HS384":
				return jwt.SigningMethodHS384
			case "HS512":
				return jwt.SigningMethodHS512
			case "RS256":
				return jwt.SigningMethodRS256
			case "RS384":
				return jwt.SigningMethodRS384
			case "RS512":
				return jwt.SigningMethodRS512
			case "ES256":
				return jwt.SigningMethodES256
			case "ES384":
				return jwt.SigningMethodES384
			case "ES512":
				return jwt.SigningMethodES512
			case "PS256":
				return jwt.SigningMethodPS256
			case "PS384":
				return jwt.SigningMethodPS384
			case "PS512":
				return jwt.SigningMethodPS512
			default:
				panic(fmt.Sprintf("JWT signing method %s is not allow", method))
			}
		}(),
		JWTSigningKey: []byte(os.Getenv("JWT_SIGNING_KEY")),
		JWTAllowMethod: func() authservice.AllowSigningMethod {
			allow := authservice.AllowSigningMethod{}
			for _, a := range strings.Split(os.Getenv("JWT_ALLOW_METHOD"), ",") {
				a = strings.TrimSpace(a)
				allow.Allowed(a)
			}
			return allow
		}(),
		JWTExpiredIn: func() time.Duration {
			var (
				t   int
				err error
			)

			if t, err = strconv.Atoi(os.Getenv("JWT_EXPIRED_IN")); err != nil {
				t = 30
			}

			return time.Duration(t * int(time.Second))
		}(),
	}

	log = logrus.WithField("app_name", options.AppName)
}

func waitForShutdownSignal() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until signaled
	sig := <-c

	return sig.String()
}

func shutdownOnSignal(svc *internal.Service) {
	signalName := waitForShutdownSignal()
	log.WithField("signal", signalName).Info("Received signal, starting shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if svc.Shutdown(ctx) {
		log.Info("Shutdown complete")
	} else {
		log.Info("Shutdown timed out")
	}
}
