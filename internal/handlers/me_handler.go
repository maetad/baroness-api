package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

type MeHandler struct {
	log *logrus.Entry
}

func NewMeHandler(log *logrus.Entry) *MeHandler {
	return &MeHandler{log}
}

func (h *MeHandler) Get(c *gin.Context) {
	var (
		user *userservice.User
		ok   bool
	)

	if user, ok = c.MustGet("user").(*userservice.User); !ok {
		h.log.Error(`Delete(): c.MustGet("user") is not *userservice.User`)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)
}
