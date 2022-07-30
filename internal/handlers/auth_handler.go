package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/config"
	"github.com/pakkaparn/no-idea-api/internal/services/authservice"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	log         *logrus.Entry
	options     config.Options
	authservice authservice.AuthServiceInterface
	userservice userservice.UserServiceInterface
}

func NewAuthHandler(
	log *logrus.Entry,
	options config.Options,
	authservice authservice.AuthServiceInterface,
	userservice userservice.UserServiceInterface,
) *AuthHandler {
	return &AuthHandler{log, options, authservice, userservice}
}

func (h *AuthHandler) Login(c *gin.Context) {
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

	token, err := h.authservice.GenerateToken(user.(*userservice.User), h.options.JWTExpiredIn)
	if err != nil {
		h.log.WithError(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Authorize(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")
	if _, err := h.authservice.ParseToken(token); err != nil {
		h.log.WithError(err).Errorf("Authorize(): h.authservice.ParseToken error %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
