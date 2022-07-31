package fieldservice

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/maetad/baroness-api/internal/model"
)

type FieldCreateRequest struct {
	WorkflowID uint            `json:"-"`
	Name       string          `json:"name" binding:"required"`
	Type       model.FieldType `json:"type" binding:"required,enum_fieldtype"`
}

type FieldUpdateRequest struct {
	Name string          `json:"name" binding:"required"`
	Type model.FieldType `json:"type" binding:"required,enum_fieldtype"`
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("enum_fieldtype", func(fl validator.FieldLevel) bool {
			v := fl.Field().Interface().(model.FieldType)
			return v.IsValid()
		})
	}
}
