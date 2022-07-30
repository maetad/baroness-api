package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
		h.log.WithError(err).Errorf("Login(): h.userservice.GetByUsername error %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = user.ValidatePassword(req.Password); err != nil {
		h.log.WithError(err).Errorf("Login(): user.ValidatePassword error %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := h.authservice.GenerateToken(user.(*userservice.User), h.options.JWTExpiredIn)
	if err != nil {
		h.log.WithError(err).Errorf("Login(): h.authservice.GenerateToken error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Authorize(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")

	var (
		claims jwt.MapClaims
		user   userservice.UserInterface
		err    error
	)

	if claims, err = h.authservice.ParseToken(token); err != nil {
		h.log.WithError(err).Errorf("Authorize(): h.authservice.ParseToken error %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims["username"] == nil {
		h.log.Error("Authorize(): claims username not exists")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if _, ok := claims["username"].(string); !ok {
		h.log.Error("Authorize(): claims usernamed is not string")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if user, err = h.userservice.GetByUsername(claims["username"].(string)); err != nil {
		h.log.WithError(err).Errorf("Authorize(): h.userservice.Get error %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)

	c.Next()
}
