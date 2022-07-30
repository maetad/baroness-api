package model

type State struct {
	Model
	WorkflowID uint     `json:"workflow_id"`
	Workflow   Workflow `json:"event" gorm:"foreignkey:WorkflowID"`
	Name       string   `json:"name"`
	Author
}
