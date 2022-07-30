package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/fieldservice"
	"github.com/sirupsen/logrus"
)

type FieldHandler struct {
	log          *logrus.Entry
	fieldservice fieldservice.FieldServiceInterface
}

func NewFieldHandler(log *logrus.Entry, fieldservice fieldservice.FieldServiceInterface) *FieldHandler {
	return &FieldHandler{log, fieldservice}
}

func (h *FieldHandler) List(c *gin.Context) {
	workflow := c.MustGet("workflow").(*model.Workflow)

	list, err := h.fieldservice.List(workflow.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *FieldHandler) Get(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("field_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	field, err := h.fieldservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Get(): h.fieldservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Set("field", field)
	c.Next()
}

func (h *FieldHandler) Create(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	workflow := c.MustGet("workflow").(*model.Workflow)

	var r fieldservice.FieldCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	r.WorkflowID = workflow.ID

	field, err := h.fieldservice.Create(r, user)
	if err != nil {
		h.log.WithError(err).Errorf("Create(): h.fieldservice.Create error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, field)
}

func (h *FieldHandler) Show(c *gin.Context) {
	field := c.MustGet("field").(*model.Field)

	c.JSON(http.StatusOK, field)
}

func (h *FieldHandler) Update(c *gin.Context) {
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(c.Param("field_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var r fieldservice.FieldUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	field, err := h.fieldservice.Get(uint(id))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.fieldservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	field, err = h.fieldservice.Update(field, r, c.MustGet("user").(*model.User))
	if err != nil {
		h.log.WithError(err).Errorf("Update(): h.fieldservice.Update error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, field)
}

func (h *FieldHandler) Delete(c *gin.Context) {
	var (
		id    int
		err   error
		field *model.Field
	)

	if id, err = strconv.Atoi(c.Param("field_id")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if field, err = h.fieldservice.Get(uint(id)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.fieldservice.Get error %v", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = h.fieldservice.Delete(field, c.MustGet("user").(*model.User)); err != nil {
		h.log.WithError(err).Errorf("Delete(): h.fieldservice.Delete error %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
