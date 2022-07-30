package eventservice

import (
	"time"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"gorm.io/gorm"
)

type EventService struct {
	db database.DatabaseInterface
}

type EventServiceInterface interface {
	List() ([]model.Event, error)
	Create(r EventCreateRequest, creator *model.User) (*model.Event, error)
	Get(id uint) (*model.Event, error)
	Update(event *model.Event, r EventUpdateRequest, updator *model.User) (*model.Event, error)
	Delete(event *model.Event, deletor *model.User) error
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

func (s EventService) Create(r EventCreateRequest, creator *model.User) (*model.Event, error) {
	event := &model.Event{
		Name:     r.Name,
		Platform: r.Platform,
		Channel:  r.Channel,
		StartAt:  r.StartAt,
		EndAt:    r.EndAt,
		Author: model.Author{
			CreatedBy: creator.ID,
			UpdatedBy: creator.ID,
		},
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

func (s EventService) Update(event *model.Event, r EventUpdateRequest, updator *model.User) (*model.Event, error) {
	event.Name = r.Name
	event.Platform = r.Platform
	event.Channel = r.Channel
	event.StartAt = r.StartAt
	event.EndAt = r.EndAt
	event.UpdatedBy = updator.ID

	if result := s.db.Save(event); result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func (s EventService) Delete(event *model.Event, deletor *model.User) error {
	event.DeletedAt = gorm.DeletedAt{
		Time: time.Now(),
	}
	event.DeletedBy = deletor.ID

	return s.db.Save(event).Error
}
