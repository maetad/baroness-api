package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/services/authservice"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	log         *logrus.Entry
	authservice authservice.AuthServiceInterface
	userservice userservice.UserServiceInterface
}

func NewAuthHandler(
	log *logrus.Entry,
	authservice authservice.AuthServiceInterface,
	userservice userservice.UserServiceInterface,
) *AuthHandler {
	return &AuthHandler{log, authservice, userservice}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	var (
		user userservice.UserInterface
		err  error
	)

	if user, err = h.userservice.GetByUsername(req.Username); err != nil {
		h.log.WithError(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = user.ValidatePassword(req.Password); err != nil {
		h.log.WithError(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := h.authservice.GenerateToken(user.(*userservice.User))
	if err != nil {
		h.log.WithError(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
