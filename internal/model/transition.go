package model

type Transition struct {
	Model
	WorkflowID uint     `json:"workflow_id"`
	ParentID   uint     `json:"parent_id"`
	TargetID   uint     `json:"target_id"`
	Workflow   Workflow `json:"workflow" gorm:"foreignkey:WorkflowID"`
	Parent     State    `json:"parent" gorm:"foreignkey:ParentID"`
	Target     State    `json:"target" gorm:"foreignkey:TargetID"`
	Name       string   `json:"name"`
	Author
}
