package internal

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/config"
	"github.com/pakkaparn/no-idea-api/internal/database"
	"github.com/pakkaparn/no-idea-api/internal/services/authservice"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

type Service struct {
	Http *http.Server
	log  *logrus.Entry
}

type internalService struct {
	authservice authservice.AuthServiceInterface
	userservice userservice.UserServiceInterface
}

func New(
	ctx context.Context,
	l *logrus.Entry,
	options config.Options,
) (*Service, error) {
	log = l

	r := gin.Default()

	db, err := database.Connect(options)
	if err != nil {
		log.WithError(err).Fatal("database.Connect()")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatal("db.DB()")
	}

	if err = database.AutoMigration(sqlDB); err != nil {
		log.WithError(err).Fatal("database.AutoMigration()")
	}

	svc := Service{
		Http: &http.Server{
			Addr:    options.ListenAddressHTTP,
			Handler: r.Handler(),
		},
		log: l,
	}

	services := internalService{
		authservice: authservice.New(options.JWTSigningMethod, options.JWTSigningKey, options.JWTAllowMethod),
		userservice: userservice.New(db),
	}

	registerRouter(r, l, options, services)

	return &svc, nil
}

func (s *Service) Close() {
	s.Http.Close()
}

func (s *Service) Shutdown(ctx context.Context) bool {
	err := s.Http.Shutdown(ctx)

	return err == nil
}
