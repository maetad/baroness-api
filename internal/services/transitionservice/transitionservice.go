package transitionservice

import (
	"time"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"gorm.io/gorm"
)

type TransitionService struct {
	db database.DatabaseInterface
}

type TransitionServiceInterface interface {
	List(workflowID uint) ([]model.Transition, error)
	Create(r TransitionCreateRequest, creator *model.User) (*model.Transition, error)
	Get(id uint) (*model.Transition, error)
	Update(transition *model.Transition, r TransitionUpdateRequest, updator *model.User) (*model.Transition, error)
	Delete(transition *model.Transition, deletor *model.User) error
}

func New(db database.DatabaseInterface) TransitionServiceInterface {
	return TransitionService{db}
}

func (s TransitionService) List(workflowID uint) ([]model.Transition, error) {
	transitions := []model.Transition{}

	if result := s.db.Find(&transitions, "workflow_id = ?", workflowID); result.Error != nil {
		return nil, result.Error
	}

	return transitions, nil
}

func (s TransitionService) Create(r TransitionCreateRequest, creator *model.User) (*model.Transition, error) {
	transition := &model.Transition{
		Name:     r.Name,
		ParentID: r.ParentID,
		TargetID: r.TargetID,
		Author: model.Author{
			CreatedBy: creator.ID,
			UpdatedBy: creator.ID,
		},
	}

	if result := s.db.Create(transition); result.Error != nil {
		return nil, result.Error
	}

	return transition, nil
}

func (s TransitionService) Get(id uint) (*model.Transition, error) {
	transition := &model.Transition{}

	if result := s.db.First(transition, id); result.Error != nil {
		return nil, result.Error
	}

	return transition, nil
}

func (s TransitionService) Update(transition *model.Transition, r TransitionUpdateRequest, updator *model.User) (*model.Transition, error) {
	transition.Name = r.Name
	transition.ParentID = r.ParentID
	transition.TargetID = r.TargetID
	transition.UpdatedBy = updator.ID

	if result := s.db.Save(transition); result.Error != nil {
		return nil, result.Error
	}

	return transition, nil
}

func (s TransitionService) Delete(transition *model.Transition, deletor *model.User) error {
	transition.DeletedAt = gorm.DeletedAt{
		Time: time.Now(),
	}
	transition.DeletedBy = deletor.ID

	return s.db.Save(transition).Error
}
