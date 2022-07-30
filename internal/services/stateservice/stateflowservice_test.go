package stateservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/stateservice"
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
			name: "new state",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(stateservice.StateService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(stateservice.StateService{}).Kind())
			}
		})
	}
}

func TestStateService_List(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		workflowID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.State
		wantErr bool
	}{
		{
			name: "listed success",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything, "workflow_id = ?", mock.Anything).
					Return(&gorm.DB{})

				return fields{db}
			}(),
			args: args{
				workflowID: 1,
			},
			want: []model.State{},
		},
		{
			name: "listed fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything, "workflow_id = ?", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("find error"),
					})

				return fields{db}
			}(),
			args: args{
				workflowID: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stateservice.New(tt.fields.db)
			got, err := s.List(tt.args.workflowID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Create(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r       stateservice.StateCreateRequest
		creator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.State
		wantErr bool
	}{
		{
			name: "state created",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				r: stateservice.StateCreateRequest{
					WorkflowID: 1,
					Name:       "state name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.State{
				WorkflowID: 1,
				Name:       "state name",
				Author: model.Author{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "state create fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("error"),
					})

				return fields{db}
			}(),
			args: args{
				r: stateservice.StateCreateRequest{
					WorkflowID: 1,
					Name:       "state name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stateservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Get(t *testing.T) {
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
		want    *model.State
		wantErr bool
	}{
		{
			name: "state found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.State"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &model.State{},
		},
		{
			name: "state not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.State"), uint(1)).
					Return(&gorm.DB{
						Error: errors.New("state not found"),
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
			s := stateservice.New(tt.fields.db)
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Update(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		state   *model.State
		r       stateservice.StateUpdateRequest
		updator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.State
		wantErr bool
	}{
		{
			name: "state updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.State")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				state: &model.State{
					Name: "state name",
				},
				r: stateservice.StateUpdateRequest{
					Name: "new state name",
				},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.State{
				Name: "new state name",
				Author: model.Author{
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "state update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.State")).
					Return(&gorm.DB{
						Error: errors.New("state update error"),
					})

				return fields{db}
			}(),
			args: args{
				state:   &model.State{},
				r:       stateservice.StateUpdateRequest{},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stateservice.New(tt.fields.db)
			got, err := s.Update(tt.args.state, tt.args.r, tt.args.updator)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateService_Delete(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		state   *model.State
		deletor *model.User
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
				db.On("Save", mock.AnythingOfType("*model.State")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				state:   &model.State{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.State")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				state:   &model.State{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stateservice.New(tt.fields.db)
			if err := s.Delete(tt.args.state, tt.args.deletor); (err != nil) != tt.wantErr {
				t.Errorf("StateService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
