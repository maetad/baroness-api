package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	log         *logrus.Entry
	userservice userservice.UserServiceInterface
}

func NewUserHandler(log *logrus.Entry, userservice userservice.UserServiceInterface) *UserHandler {
	return &UserHandler{log, userservice}
}

func (h *UserHandler) List(c *gin.Context) {
	list, err := h.userservice.List()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}
