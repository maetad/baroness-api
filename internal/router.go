package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/config"
	"github.com/maetad/baroness-api/internal/handlers"
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

	r.POST("/auth/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(authHandler.Authorize)
	{
		meHandler := handlers.NewMeHandler(l, services.userservice)
		authorized.GET("/me", meHandler.Get)
		authorized.PUT("/me", meHandler.Update)

		userRoute := authorized.Group("/users")
		{
			userHandler := handlers.NewUserHandler(l, services.userservice)
			userRoute.GET("/", userHandler.List)
			userRoute.POST("/", userHandler.Create)
			userRoute.GET("/:id", userHandler.Get)
			userRoute.PUT("/:id", userHandler.Update)
			userRoute.DELETE("/:id", userHandler.Delete)
		}

		eventRoute := authorized.Group("/events")
		{
			eventHandler := handlers.NewEventHandler(l, services.eventservice)
			eventRoute.GET("/", eventHandler.List)
			eventRoute.POST("/", eventHandler.Create)
			eventRoute.GET("/:id", eventHandler.Get)
			eventRoute.PUT("/:id", eventHandler.Update)
			eventRoute.DELETE("/:id", eventHandler.Delete)
		}
	}
}
