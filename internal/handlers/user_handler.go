package handlers

import (
	"net/http"
	"strconv"

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

func (h *UserHandler) Create(c *gin.Context) {
	var r userservice.UserCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	user, err := h.userservice.Create(r)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.userservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user, err := h.userservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.userservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r userservice.UserUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	user, err := h.userservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.userservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user, err = h.userservice.Update(user, r)
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.userservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user, err := h.userservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Delete(): h.userservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.userservice.Delete(user); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.userservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
