package model

import (
	"time"

	"github.com/lib/pq"
)

type Event struct {
	Model
	Name     string         `json:"name"`
	Platform pq.StringArray `json:"platform" gorm:"type:event_platform[]"`
	Channel  pq.StringArray `json:"channel" gorm:"type:event_channel[]"`
	StartAt  time.Time      `json:"start_at"`
	EndAt    time.Time      `json:"end_at"`
	Author
}