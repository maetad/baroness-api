package model

type Workflow struct {
	Model
	EventID uint   `json:"event_id"`
	Event   Event  `json:"event" gorm:"foreignkey:EventID"`
	Name    string `json:"name"`
	Author
}
