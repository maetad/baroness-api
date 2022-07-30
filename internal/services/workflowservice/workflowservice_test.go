package workflowservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/workflowservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var db = &mocks.DatabaseInterface{}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "new workflow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workflowservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(workflowservice.WorkflowService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(workflowservice.WorkflowService{}).Kind())
			}
		})
	}
}

func TestWorkflowService_List(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		eventID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Workflow
		wantErr bool
	}{
		{
			name: "listed success",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything, "event_id = ?", mock.Anything).
					Return(&gorm.DB{})

				return fields{db}
			}(),
			args: args{
				eventID: 1,
			},
			want: []model.Workflow{},
		},
		{
			name: "listed fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything, "event_id = ?", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("find error"),
					})

				return fields{db}
			}(),
			args: args{
				eventID: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := workflowservice.New(tt.fields.db)
			got, err := s.List(tt.args.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkflowService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkflowService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowService_Create(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r       workflowservice.WorkflowCreateRequest
		creator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Workflow
		wantErr bool
	}{
		{
			name: "workflow created",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				r: workflowservice.WorkflowCreateRequest{
					EventID: 1,
					Name:    "workflow name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Workflow{
				EventID: 1,
				Name:    "workflow name",
				Author: model.Author{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "workflow create fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("error"),
					})

				return fields{db}
			}(),
			args: args{
				r: workflowservice.WorkflowCreateRequest{
					EventID: 1,
					Name:    "workflow name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := workflowservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkflowService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkflowService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowService_Get(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Workflow
		wantErr bool
	}{
		{
			name: "workflow found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Workflow"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &model.Workflow{},
		},
		{
			name: "workflow not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Workflow"), uint(1)).
					Return(&gorm.DB{
						Error: errors.New("workflow not found"),
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := workflowservice.New(tt.fields.db)
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkflowService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkflowService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowService_Update(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		workflow *model.Workflow
		r        workflowservice.WorkflowUpdateRequest
		updator  *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Workflow
		wantErr bool
	}{
		{
			name: "workflow updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Workflow")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				workflow: &model.Workflow{
					Name: "workflow name",
				},
				r: workflowservice.WorkflowUpdateRequest{
					Name: "new workflow name",
				},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Workflow{
				Name: "new workflow name",
				Author: model.Author{
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "workflow update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Workflow")).
					Return(&gorm.DB{
						Error: errors.New("workflow update error"),
					})

				return fields{db}
			}(),
			args: args{
				workflow: &model.Workflow{},
				r:        workflowservice.WorkflowUpdateRequest{},
				updator:  &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := workflowservice.New(tt.fields.db)
			got, err := s.Update(tt.args.workflow, tt.args.r, tt.args.updator)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkflowService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkflowService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowService_Delete(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		workflow *model.Workflow
		deletor  *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete success",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Workflow")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				workflow: &model.Workflow{},
				deletor:  &model.User{Model: model.Model{ID: 1}},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Workflow")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				workflow: &model.Workflow{},
				deletor:  &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := workflowservice.New(tt.fields.db)
			if err := s.Delete(tt.args.workflow, tt.args.deletor); (err != nil) != tt.wantErr {
				t.Errorf("WorkflowService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
