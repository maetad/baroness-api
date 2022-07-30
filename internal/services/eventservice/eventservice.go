package eventservice

import (
	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
)

type EventService struct {
	db database.DatabaseInterface
}

type EventServiceInterface interface {
	List() ([]model.Event, error)
	Create(r EventCreateRequest) (*model.Event, error)
	Get(id uint) (*model.Event, error)
	Update(event *model.Event, r EventUpdateRequest) (*model.Event, error)
	Delete(event *model.Event) error
}

func New(db database.DatabaseInterface) EventServiceInterface {
	return EventService{db}
}

func (s EventService) List() ([]model.Event, error) {
	events := []model.Event{}
	if result := s.db.Find(&events); result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func (s EventService) Create(r EventCreateRequest) (*model.Event, error) {
	event := &model.Event{
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

func (s EventService) Get(id uint) (*model.Event, error) {
	event := &model.Event{}

	if result := s.db.First(event, id); result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func (s EventService) Update(event *model.Event, r EventUpdateRequest) (*model.Event, error) {
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

func (s EventService) Delete(event *model.Event) error {
	return s.db.Delete(event).Error
}
