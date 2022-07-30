package model

import "github.com/lib/pq"

type FieldType string

const (
	FieldTypeText      FieldType = "text"
	FieldTypeParagraph FieldType = "paragraph"
	FieldTypeDate      FieldType = "date"
	FieldTypeNumber    FieldType = "number"
	FieldTypeTimestamp FieldType = "timestamp"
	FieldTypeDropdown  FieldType = "dropdown"
	FieldTypeCheckbox  FieldType = "checkbox"
)

type Field struct {
	Model
	WorkflowID uint           `json:"workflow_id"`
	Workflow   Workflow       `json:"workflow" gorm:"foreignkey:WorkflowID"`
	Name       string         `json:"name"`
	Type       FieldType      `json:"type"`
	Options    pq.StringArray `json:"options"`
	Author
}
