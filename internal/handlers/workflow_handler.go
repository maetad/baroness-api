package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/workflowservice"
	"github.com/sirupsen/logrus"
)

type WorkflowHandler struct {
	log             *logrus.Entry
	workflowservice workflowservice.WorkflowServiceInterface
}

func NewWorkflowHandler(log *logrus.Entry, workflowservice workflowservice.WorkflowServiceInterface) *WorkflowHandler {
	return &WorkflowHandler{log, workflowservice}
}

func (h *WorkflowHandler) List(c *gin.Context) {
	event := c.MustGet("event").(*model.Event)

	list, err := h.workflowservice.List(event.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *WorkflowHandler) Create(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	event := c.MustGet("event").(*model.Event)

	var r workflowservice.WorkflowCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	r.EventID = event.ID

	workflow, err := h.workflowservice.Create(r, user)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.workflowservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

func (h *WorkflowHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("workflow_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	workflow, err := h.workflowservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.workflowservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("workflow_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r workflowservice.WorkflowUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	workflow, err := h.workflowservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.workflowservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	workflow, err = h.workflowservice.Update(workflow, r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.workflowservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) Delete(c *gin.Context) {
	var (
		id       int
		err      error
		workflow *model.Workflow
	)

	if id, err = strconv.Atoi(c.Param("workflow_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if workflow, err = h.workflowservice.Get(uint(id)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.workflowservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.workflowservice.Delete(workflow, c.MustGet("user").(*model.User)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.workflowservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
