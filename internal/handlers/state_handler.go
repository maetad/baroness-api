package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/stateservice"
	"github.com/sirupsen/logrus"
)

type StateHandler struct {
	log          *logrus.Entry
	stateservice stateservice.StateServiceInterface
}

func NewStateHandler(log *logrus.Entry, stateservice stateservice.StateServiceInterface) *StateHandler {
	return &StateHandler{log, stateservice}
}

func (h *StateHandler) List(c *gin.Context) {
	workflow := c.MustGet("workflow").(*model.Workflow)

	list, err := h.stateservice.List(workflow.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *StateHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("state_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	state, err := h.stateservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.stateservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Set("state", state)
	c.Next()
}

func (h *StateHandler) Create(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	workflow := c.MustGet("workflow").(*model.Workflow)

	var r stateservice.StateCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	r.WorkflowID = workflow.ID

	state, err := h.stateservice.Create(r, user)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.stateservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, state)
}

func (h *StateHandler) Show(c *gin.Context) {
	state := c.MustGet("state").(*model.State)

	c.JSON(http.StatusOK, state)
}

func (h *StateHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("state_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r stateservice.StateUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	state, err := h.stateservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.stateservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	state, err = h.stateservice.Update(state, r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.stateservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, state)
}

func (h *StateHandler) Delete(c *gin.Context) {
	var (
		id    int
		err   error
		state *model.State
	)

	if id, err = strconv.Atoi(c.Param("state_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if state, err = h.stateservice.Get(uint(id)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.stateservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.stateservice.Delete(state, c.MustGet("user").(*model.User)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.stateservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
