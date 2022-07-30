package eventservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/eventservice"
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
			name: "new event",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eventservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(eventservice.EventService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(eventservice.EventService{}).Kind())
			}
		})
	}
}

func TestEventService_List(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Event
		wantErr bool
	}{
		{
			name: "listed success",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			want: []model.Event{},
		},
		{
			name: "listed fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Find", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("find error"),
					})

				return fields{db}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			got, err := s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("EventService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_Create(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r       eventservice.EventCreateRequest
		creator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		{
			name: "event created",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				r: eventservice.EventCreateRequest{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Event{
				Name:     "event name",
				Platform: []string{"platform 1", "platform 2"},
				Channel:  []string{"channel 1", "channel 2"},
				Author: model.Author{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "event create fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Create", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("error"),
					})

				return fields{db}
			}(),
			args: args{
				r: eventservice.EventCreateRequest{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				},
				creator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_Get(t *testing.T) {
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
		want    *model.Event
		wantErr bool
	}{
		{
			name: "event found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Event"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &model.Event{},
		},
		{
			name: "event not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*model.Event"), uint(1)).
					Return(&gorm.DB{
						Error: errors.New("event not found"),
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
			s := eventservice.New(tt.fields.db)
			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_Update(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		event   *model.Event
		r       eventservice.EventUpdateRequest
		updator *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		{
			name: "event updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Event")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				event: &model.Event{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				},
				r: eventservice.EventUpdateRequest{
					Name:     "new event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			want: &model.Event{
				Name:     "new event name",
				Platform: []string{"platform 1", "platform 2"},
				Channel:  []string{"channel 1", "channel 2"},
				Author: model.Author{
					UpdatedBy: 1,
				},
			},
		},
		{
			name: "event update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Event")).
					Return(&gorm.DB{
						Error: errors.New("event update error"),
					})

				return fields{db}
			}(),
			args: args{
				event:   &model.Event{},
				r:       eventservice.EventUpdateRequest{},
				updator: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			got, err := s.Update(tt.args.event, tt.args.r, tt.args.updator)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_Delete(t *testing.T) {
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		event   *model.Event
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
				db.On("Save", mock.AnythingOfType("*model.Event")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				event:   &model.Event{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*model.Event")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				event:   &model.Event{},
				deletor: &model.User{Model: model.Model{ID: 1}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			if err := s.Delete(tt.args.event, tt.args.deletor); (err != nil) != tt.wantErr {
				t.Errorf("EventService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
