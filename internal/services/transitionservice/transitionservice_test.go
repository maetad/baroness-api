package transitionservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/transitionservice"
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
			name: "new transition",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transitionservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(transitionservice.TransitionService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(transitionservice.TransitionService{}).Kind())
			}
		})
	}
}

func TestTransitionService_List(t *testing.T) {
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
		want    []model.Transition
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
			want: []model.Transition{},
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
			s := transitionservice.New(tt.fields.db)
			got, err := s.List(tt.args.workflowID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransitionService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionService_Create(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r       transitionservice.TransitionCreateRequest
		creator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Transition
		wantErr bool
	}{
		{
			name: "transition created",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				r: transitionservice.TransitionCreateRequest{
					Name:     "transition name",
					ParentID: 1,
					TargetID: 2,
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Transition{
				Name:     "transition name",
				ParentID: 1,
				TargetID: 2,
				Author: model.Author{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "transition create fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("error"),
					})

				return fields{db}
			}(),
			args: args{
				r: transitionservice.TransitionCreateRequest{
					Name:     "transition name",
					ParentID: 1,
					TargetID: 2,
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := transitionservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransitionService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionService_Get(t *testing.T) {
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
		want    *model.Transition
		wantErr bool
	}{
		{
			name: "transition found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Transition"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &model.Transition{},
		},
		{
			name: "transition not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Transition"), uint(1)).
					Return(&gorm.DB{
						Error: errors.New("transition not found"),
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
			s := transitionservice.New(tt.fields.db)
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransitionService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionService_Update(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		transition *model.Transition
		r          transitionservice.TransitionUpdateRequest
		updator    *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Transition
		wantErr bool
	}{
		{
			name: "transition updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Transition")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				transition: &model.Transition{
					Name: "transition name",
				},
				r: transitionservice.TransitionUpdateRequest{
					Name: "new transition name",
				},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Transition{
				Name: "new transition name",
				Author: model.Author{
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "transition update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Transition")).
					Return(&gorm.DB{
						Error: errors.New("transition update error"),
					})

				return fields{db}
			}(),
			args: args{
				transition: &model.Transition{},
				r:          transitionservice.TransitionUpdateRequest{},
				updator:    &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := transitionservice.New(tt.fields.db)
			got, err := s.Update(tt.args.transition, tt.args.r, tt.args.updator)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransitionService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionService_Delete(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		transition *model.Transition
		deletor    *model.User
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
				db.On("Save", mock.AnythingOfType("*model.Transition")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				transition: &model.Transition{},
				deletor:    &model.User{Model: model.Model{ID: 1}},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Transition")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				transition: &model.Transition{},
				deletor:    &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := transitionservice.New(tt.fields.db)
			if err := s.Delete(tt.args.transition, tt.args.deletor); (err != nil) != tt.wantErr {
				t.Errorf("TransitionService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
