package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	event, err := h.eventservice.Create(r)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.eventservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	event, err := h.eventservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.eventservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r eventservice.EventUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	event, err := h.eventservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.eventservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	event, err = h.eventservice.Update(event, r)
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.eventservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	var (
		currentEvent *eventservice.Event
		ok           bool
		id           int
		err          error
		event        *eventservice.Event
	)

	if currentEvent, ok = c.MustGet("event").(*eventservice.Event); !ok {
		h.log.Error(`Delete(): c.MustGet("event") is not *eventservice.Event`)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if currentEvent.ID == uint(id) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if event, err = h.eventservice.Get(uint(id)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.eventservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.eventservice.Delete(event); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.eventservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
