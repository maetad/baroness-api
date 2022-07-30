package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

type MeHandler struct {
	log         *logrus.Entry
	userservice userservice.UserServiceInterface
}

func NewMeHandler(
	log *logrus.Entry,
	userservice userservice.UserServiceInterface,
) *MeHandler {
	return &MeHandler{log, userservice}
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

func (h *MeHandler) Update(c *gin.Context) {
	var (
		user *userservice.User
		r    userservice.UserUpdateRequest
		ok   bool
	)

	if user, ok = c.MustGet("user").(*userservice.User); !ok {
		h.log.Error(`Delete(): c.MustGet("user") is not *userservice.User`)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	u, err := h.userservice.Update(user, r)
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.userservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u)
}
