package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/eventservice"
	"github.com/sirupsen/logrus"
)

type EventHandler struct {
	log          *logrus.Entry
	eventservice eventservice.EventServiceInterface
}

func NewEventHandler(log *logrus.Entry, eventservice eventservice.EventServiceInterface) *EventHandler {
	return &EventHandler{log, eventservice}
}

func (h *EventHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("event_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	event, err := h.eventservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.eventservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Set("event", event)
	c.Next()
}

func (h *EventHandler) List(c *gin.Context) {
	list, err := h.eventservice.List()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *EventHandler) Create(c *gin.Context) {
	var r eventservice.EventCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	event, err := h.eventservice.Create(r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.eventservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Show(c *gin.Context) {
	event := c.MustGet("event").(*model.Event)

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Update(c *gin.Context) {
	var r eventservice.EventUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	event := c.MustGet("event").(*model.Event)

	event, err := h.eventservice.Update(event, r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.eventservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	event := c.MustGet("event").(*model.Event)

	if err := h.eventservice.Delete(event, c.MustGet("user").(*model.User)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.eventservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
