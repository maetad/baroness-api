package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/transitionservice"
	"github.com/sirupsen/logrus"
)

type TransitionHandler struct {
	log               *logrus.Entry
	transitionservice transitionservice.TransitionServiceInterface
}

func NewTransitionHandler(log *logrus.Entry, transitionservice transitionservice.TransitionServiceInterface) *TransitionHandler {
	return &TransitionHandler{log, transitionservice}
}

func (h *TransitionHandler) List(c *gin.Context) {
	workflow := c.MustGet("workflow").(*model.Workflow)

	list, err := h.transitionservice.List(workflow.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *TransitionHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("transition_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	transition, err := h.transitionservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.transitionservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Set("transition", transition)
	c.Next()
}

func (h *TransitionHandler) Create(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	workflow := c.MustGet("workflow").(*model.Workflow)

	var r transitionservice.TransitionCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	r.WorkflowID = workflow.ID

	transition, err := h.transitionservice.Create(r, user)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.transitionservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, transition)
}

func (h *TransitionHandler) Show(c *gin.Context) {
	transition := c.MustGet("transition").(*model.Transition)

	c.JSON(http.StatusOK, transition)
}

func (h *TransitionHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("transition_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r transitionservice.TransitionUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	transition, err := h.transitionservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.transitionservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	transition, err = h.transitionservice.Update(transition, r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.transitionservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, transition)
}

func (h *TransitionHandler) Delete(c *gin.Context) {
	var (
		id         int
		err        error
		transition *model.Transition
	)

	if id, err = strconv.Atoi(c.Param("transition_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if transition, err = h.transitionservice.Get(uint(id)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.transitionservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.transitionservice.Delete(transition, c.MustGet("user").(*model.User)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.transitionservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
