package eventservice_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/pakkaparn/no-idea-api/internal/database"
	"github.com/pakkaparn/no-idea-api/internal/services/eventservice"
	"github.com/pakkaparn/no-idea-api/mocks"
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
		want    []eventservice.Event
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
			want: []eventservice.Event{},
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
	var now = time.Now()

	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		r eventservice.EventCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *eventservice.Event
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
					StartAt:  now,
					EndAt:    now,
				},
			},
			want: &eventservice.Event{
				Name:     "event name",
				Platform: []string{"platform 1", "platform 2"},
				Channel:  []string{"channel 1", "channel 2"},
				StartAt:  now,
				EndAt:    now,
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
					StartAt:  now,
					EndAt:    now,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			got, err := s.Create(tt.args.r)
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
		want    *eventservice.Event
		wantErr bool
	}{
		{
			name: "event found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*eventservice.Event"), uint(1)).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				id: 1,
			},
			want: &eventservice.Event{},
		},
		{
			name: "event not found",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("First", mock.AnythingOfType("*eventservice.Event"), uint(1)).
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
	var now = time.Now()
	type fields struct {
		db database.DatabaseInterface
	}
	type args struct {
		event *eventservice.Event
		r     eventservice.EventUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *eventservice.Event
		wantErr bool
	}{
		{
			name: "event updated",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*eventservice.Event")).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			args: args{
				event: &eventservice.Event{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
					StartAt:  now,
					EndAt:    now,
				},
				r: eventservice.EventUpdateRequest{
					Name:     "new event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
					StartAt:  now,
					EndAt:    now,
				},
			},
			want: &eventservice.Event{
				Name:     "new event name",
				Platform: []string{"platform 1", "platform 2"},
				Channel:  []string{"channel 1", "channel 2"},
				StartAt:  now,
				EndAt:    now,
			},
		},
		{
			name: "event update error",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Save", mock.AnythingOfType("*eventservice.Event")).
					Return(&gorm.DB{
						Error: errors.New("event update error"),
					})

				return fields{db}
			}(),
			args: args{
				event: &eventservice.Event{},
				r:     eventservice.EventUpdateRequest{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			got, err := s.Update(tt.args.event, tt.args.r)
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
		event *eventservice.Event
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
				db.On("Delete", mock.AnythingOfType("*eventservice.Event")).
					Return(&gorm.DB{
						Error: nil,
					})
				return fields{db}
			}(),
			args: args{
				event: &eventservice.Event{},
			},
		},
		{
			name: "delete fail",
			fields: func() fields {
				db := &mocks.DatabaseInterface{}
				db.On("Delete", mock.AnythingOfType("*eventservice.Event")).
					Return(&gorm.DB{
						Error: errors.New("delete fail"),
					})
				return fields{db}
			}(),
			args: args{
				event: &eventservice.Event{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := eventservice.New(tt.fields.db)
			if err := s.Delete(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("EventService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
