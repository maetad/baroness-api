package workflowservice

import (
	"time"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"gorm.io/gorm"
)

type WorkflowService struct {
	db database.DatabaseInterface
}

type WorkflowServiceInterface interface {
	List(eventID uint) ([]model.Workflow, error)
	Create(r WorkflowCreateRequest, creator *model.User) (*model.Workflow, error)
	Get(id uint) (*model.Workflow, error)
	Update(workflow *model.Workflow, r WorkflowUpdateRequest, updator *model.User) (*model.Workflow, error)
	Delete(workflow *model.Workflow, deletor *model.User) error
}

func New(db database.DatabaseInterface) WorkflowServiceInterface {
	return WorkflowService{db}
}

func (s WorkflowService) List(eventID uint) ([]model.Workflow, error) {
	workflows := []model.Workflow{}
	if result := s.db.Find(&workflows, "event_id = ?", eventID); result.Error != nil {
		return nil, result.Error
	}

	return workflows, nil
}

func (s WorkflowService) Create(r WorkflowCreateRequest, creator *model.User) (*model.Workflow, error) {
	workflow := &model.Workflow{
		EventID: r.EventID,
		Name:    r.Name,
		Author: model.Author{
			CreatedBy: creator.ID,
			UpdatedBy: creator.ID,
		},
	}

	if result := s.db.Create(workflow); result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (s WorkflowService) Get(id uint) (*model.Workflow, error) {
	workflow := &model.Workflow{}

	if result := s.db.First(workflow, id); result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (s WorkflowService) Update(workflow *model.Workflow, r WorkflowUpdateRequest, updator *model.User) (*model.Workflow, error) {
	workflow.Name = r.Name
	workflow.UpdatedBy = updator.ID

	if result := s.db.Save(workflow); result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (s WorkflowService) Delete(workflow *model.Workflow, deletor *model.User) error {
	workflow.DeletedAt = gorm.DeletedAt{
		Time: time.Now(),
	}
	workflow.DeletedBy = deletor.ID

	return s.db.Save(workflow).Error
}
