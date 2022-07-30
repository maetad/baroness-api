package fieldservice

import (
	"time"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"gorm.io/gorm"
)

type FieldService struct {
	db database.DatabaseInterface
}

type FieldServiceInterface interface {
	List(workflowID uint) ([]model.Field, error)
	Create(r FieldCreateRequest, creator *model.User) (*model.Field, error)
	Get(id uint) (*model.Field, error)
	Update(field *model.Field, r FieldUpdateRequest, updator *model.User) (*model.Field, error)
	Delete(field *model.Field, deletor *model.User) error
}

func New(db database.DatabaseInterface) FieldServiceInterface {
	return FieldService{db}
}

func (s FieldService) List(workflowID uint) ([]model.Field, error) {
	fields := []model.Field{}
	if result := s.db.Find(&fields, "workflow_id = ?", workflowID); result.Error != nil {
		return nil, result.Error
	}

	return fields, nil
}

func (s FieldService) Create(r FieldCreateRequest, creator *model.User) (*model.Field, error) {
	field := &model.Field{
		WorkflowID: r.WorkflowID,
		Name:       r.Name,
		Author: model.Author{
			CreatedBy: creator.ID,
			UpdatedBy: creator.ID,
		},
	}

	if result := s.db.Create(field); result.Error != nil {
		return nil, result.Error
	}

	return field, nil
}

func (s FieldService) Get(id uint) (*model.Field, error) {
	field := &model.Field{}

	if result := s.db.First(field, id); result.Error != nil {
		return nil, result.Error
	}

	return field, nil
}

func (s FieldService) Update(field *model.Field, r FieldUpdateRequest, updator *model.User) (*model.Field, error) {
	field.Name = r.Name
	field.UpdatedBy = updator.ID

	if result := s.db.Save(field); result.Error != nil {
		return nil, result.Error
	}

	return field, nil
}

func (s FieldService) Delete(field *model.Field, deletor *model.User) error {
	field.DeletedAt = gorm.DeletedAt{
		Time: time.Now(),
	}
	field.DeletedBy = deletor.ID

	return s.db.Save(field).Error
}
