package fieldservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/fieldservice"
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
			name: "new field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fieldservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(fieldservice.FieldService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(fieldservice.FieldService{}).Kind())
			}
		})
	}
}

func TestFieldService_List(t *testing.T) {
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
		want    []model.Field
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
			want: []model.Field{},
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
			s := fieldservice.New(tt.fields.db)
			got, err := s.List(tt.args.workflowID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FieldService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldService_Create(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r       fieldservice.FieldCreateRequest
		creator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Field
		wantErr bool
	}{
		{
			name: "field created",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				r: fieldservice.FieldCreateRequest{
					WorkflowID: 1,
					Name:       "field name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Field{
				WorkflowID: 1,
				Name:       "field name",
				Author: model.Author{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "field create fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("error"),
					})

				return fields{db}
			}(),
			args: args{
				r: fieldservice.FieldCreateRequest{
					WorkflowID: 1,
					Name:       "field name",
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fieldservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("FieldService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldService_Get(t *testing.T) {
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
		want    *model.Field
		wantErr bool
	}{
		{
			name: "field found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Field"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &model.Field{},
		},
		{
			name: "field not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Field"), uint(1)).
					Return(&gorm.DB{
						Error: errors.New("field not found"),
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
			s := fieldservice.New(tt.fields.db)
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FieldService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldService_Update(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		field   *model.Field
		r       fieldservice.FieldUpdateRequest
		updator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Field
		wantErr bool
	}{
		{
			name: "field updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Field")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				field: &model.Field{
					Name: "field name",
				},
				r: fieldservice.FieldUpdateRequest{
					Name: "new field name",
				},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Field{
				Name: "new field name",
				Author: model.Author{
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "field update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Field")).
					Return(&gorm.DB{
						Error: errors.New("field update error"),
					})

				return fields{db}
			}(),
			args: args{
				field:   &model.Field{},
				r:       fieldservice.FieldUpdateRequest{},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fieldservice.New(tt.fields.db)
			got, err := s.Update(tt.args.field, tt.args.r, tt.args.updator)
			if (err != nil) != tt.wantErr {
				t.Errorf("FieldService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldService_Delete(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		field   *model.Field
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
				db.On("Save", mock.AnythingOfType("*model.Field")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				field:   &model.Field{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Field")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				field:   &model.Field{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fieldservice.New(tt.fields.db)
			if err := s.Delete(tt.args.field, tt.args.deletor); (err != nil) != tt.wantErr {
				t.Errorf("FieldService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
