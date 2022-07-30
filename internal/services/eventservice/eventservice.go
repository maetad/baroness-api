package eventservice

import (
	"github.com/pakkaparn/no-idea-api/internal/database"
)

type EventService struct {
	db database.DatabaseInterface
}

type EventServiceInterface interface {
	List() ([]Event, error)
	Create(r EventCreateRequest) (*Event, error)
	Get(id uint) (*Event, error)
	Update(event *Event, r EventUpdateRequest) (*Event, error)
	Delete(event *Event) error
}

func New(db database.DatabaseInterface) EventServiceInterface {
	return EventService{db}
}

func (s EventService) List() ([]Event, error) {
	events := []Event{}
	if result := s.db.Find(&events); result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func (s EventService) Create(r EventCreateRequest) (*Event, error) {
	event := &Event{
		Name:     r.Name,
		Platform: r.Platform,
		Channel:  r.Channel,
		StartAt:  r.StartAt,
		EndAt:    r.EndAt,
	}

	if result := s.db.Create(event); result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func (s EventService) Get(id uint) (*Event, error) {
	event := &Event{}

	if result := s.db.First(event, id); result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func (s EventService) Update(event *Event, r EventUpdateRequest) (*Event, error) {
	event.Name = r.Name
	event.Platform = r.Platform
	event.Channel = r.Channel
	event.StartAt = r.StartAt
	event.EndAt = r.EndAt

	if result := s.db.Save(event); result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func (s EventService) Delete(event *Event) error {
	return s.db.Delete(event).Error
}
