package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry
var timeout = 10 * time.Second

type Service struct {
	Http *http.Server
	log  *logrus.Entry
}

func New(ctx context.Context, l *logrus.Entry, options Options, r *gin.Engine) (*Service, error) {
	log = l
	svc := Service{
		Http: &http.Server{
			Addr:    options.ListenAddressHTTP,
			Handler: r.Handler(),
		},
		log: l,
	}

	registerRouter(r)

	return &svc, nil
}

func (s *Service) Close() {
	s.Http.Close()
}

func (s *Service) Shutdown(ctx context.Context) bool {
	err := s.Http.Shutdown(ctx)

	return err == nil
}
