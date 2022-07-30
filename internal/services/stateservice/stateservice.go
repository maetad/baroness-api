package stateservice

import (
	"time"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"gorm.io/gorm"
)

type StateService struct {
	db database.DatabaseInterface
}

type StateServiceInterface interface {
	List(workflowID uint) ([]model.State, error)
	Create(r StateCreateRequest, creator *model.User) (*model.State, error)
	Get(id uint) (*model.State, error)
	Update(state *model.State, r StateUpdateRequest, updator *model.User) (*model.State, error)
	Delete(state *model.State, deletor *model.User) error
}

func New(db database.DatabaseInterface) StateServiceInterface {
	return StateService{db}
}

func (s StateService) List(workflowID uint) ([]model.State, error) {
	states := []model.State{}
	if result := s.db.Find(&states, "workflow_id = ?", workflowID); result.Error != nil {
		return nil, result.Error
	}

	return states, nil
}

func (s StateService) Create(r StateCreateRequest, creator *model.User) (*model.State, error) {
	state := &model.State{
		WorkflowID: r.WorkflowID,
		Name:       r.Name,
		Author: model.Author{
			CreatedBy: creator.ID,
			UpdatedBy: creator.ID,
		},
	}

	if result := s.db.Create(state); result.Error != nil {
		return nil, result.Error
	}

	return state, nil
}

func (s StateService) Get(id uint) (*model.State, error) {
	state := &model.State{}

	if result := s.db.First(state, id); result.Error != nil {
		return nil, result.Error
	}

	return state, nil
}

func (s StateService) Update(state *model.State, r StateUpdateRequest, updator *model.User) (*model.State, error) {
	state.Name = r.Name
	state.UpdatedBy = updator.ID

	if result := s.db.Save(state); result.Error != nil {
		return nil, result.Error
	}

	return state, nil
}

func (s StateService) Delete(state *model.State, deletor *model.User) error {
	state.DeletedAt = gorm.DeletedAt{
		Time: time.Now(),
	}
	state.DeletedBy = deletor.ID

	return s.db.Save(state).Error
}
