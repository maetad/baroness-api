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

	authHandler := handlers.NewAuthHandler(l, o, services.authservice, services.userservice)

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", authHandler.Login)
	}

	userRoute := r.Group("/users").Use(authHandler.Authorize)
	{
		userHandler := handlers.NewUserHandler(l, services.userservice)
		userRoute.GET("/", userHandler.List)
		userRoute.POST("/", userHandler.Create)
	}
}
