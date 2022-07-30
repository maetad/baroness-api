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

		userHandler := handlers.NewUserHandler(l, services.userservice)
		userRoute := authorized.Group("/users")
		{
			userRoute.GET("/", userHandler.List)
			userRoute.POST("/", userHandler.Create)
			userRoute.GET("/:id", userHandler.Get)
			userRoute.PUT("/:id", userHandler.Update)
			userRoute.DELETE("/:id", userHandler.Delete)
		}

		eventHandler := handlers.NewEventHandler(l, services.eventservice)
		eventRoute := authorized.Group("/events")
		{
			eventRoute.GET("/", eventHandler.List)
			eventRoute.POST("/", eventHandler.Create)
			eventRoute.GET("/:event_id", eventHandler.Get, eventHandler.Show)
			eventRoute.PUT("/:event_id", eventHandler.Get, eventHandler.Update)
			eventRoute.DELETE("/:event_id", eventHandler.Get, eventHandler.Delete)
		}

		workflowHandler := handlers.NewWorkflowHandler(l, services.workflowservice)
		workflowRoute := authorized.Group("/events/:event_id/workflows")
		workflowRoute.Use(eventHandler.Get)
		{
			workflowRoute.GET("/", workflowHandler.List)
			workflowRoute.POST("/", workflowHandler.Create)
			workflowRoute.GET("/:workflow_id", workflowHandler.Get, workflowHandler.Show)
			workflowRoute.PUT("/:workflow_id", workflowHandler.Get, workflowHandler.Update)
			workflowRoute.DELETE("/:workflow_id", workflowHandler.Get, workflowHandler.Delete)
		}

		stateHandler := handlers.NewStateHandler(l, services.stateservice)
		stateRoute := authorized.Group("/events/:event_id/workflows/:workflow_id/states")
		stateRoute.Use(workflowHandler.Get)
		{
			stateRoute.GET("/", stateHandler.List)
			stateRoute.POST("/", stateHandler.Create)
			stateRoute.GET("/:state_id", stateHandler.Get, stateHandler.Show)
			stateRoute.PUT("/:state_id", stateHandler.Get, stateHandler.Update)
			stateRoute.DELETE("/:state_id", stateHandler.Get, stateHandler.Delete)
		}
	}
}
