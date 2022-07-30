package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/config"
	"github.com/pakkaparn/no-idea-api/internal/handlers"
	"github.com/sirupsen/logrus"
)

func registerRouter(
	r *gin.Engine,
	l *logrus.Entry,
	o config.Options,
	services internalService,
) {
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	authRoute := r.Group("/auth")
	{
		authHandler := handlers.NewAuthHandler(l, o, services.authservice, services.userservice)
		authRoute.POST("/login", authHandler.LoginHandler)
	}
}
