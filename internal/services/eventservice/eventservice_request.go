package eventservice

import "time"

type EventCreateRequest struct {
	Name     string    `json:"name" binding:"required"`
	Platform []string  `json:"platform" binding:"required"`
	Channel  []string  `json:"channel" binding:"required"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

type EventUpdateRequest struct {
	Name     string    `json:"name" binding:"required"`
	Platform []string  `json:"platform" binding:"required"`
	Channel  []string  `json:"channel" binding:"required"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}
