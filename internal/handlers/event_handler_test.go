package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maetad/baroness-api/internal/handlers"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/eventservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
)

func TestNewEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.EventHandler
	}{
		{
			name: "create event handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewEventHandler(tt.args.log, tt.args.eventservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.EventHandler{}) {
				t.Errorf("NewEventHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "listed success",
			fields: func() fields {
				events := make([]eventservice.Event, 1)
				eventservice := &mocks.EventServiceInterface{}
				eventservice.On("List").
					Return(events, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: eventservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "listed fail",
			fields: func() fields {
				eventservice := &mocks.EventServiceInterface{}
				eventservice.On("List").
					Return(nil, errors.New("list fail"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: eventservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewEventHandler(tt.fields.log, tt.fields.eventservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestEventHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "event created success",
			fields: func() fields {
				u := &mocks.EventServiceInterface{}
				u.On("Create", eventservice.EventCreateRequest{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				}).Return(&eventservice.Event{}, nil)
				return fields{
					eventservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"event name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "event created fail invalid payload",
			fields: func() fields {
				u := &mocks.EventServiceInterface{}
				u.On("Create", eventservice.EventCreateRequest{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				}).Return(&eventservice.Event{}, nil)
				return fields{
					eventservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"event name","platform":"platform 1","channel":["channel 1", "channel 2"]}`)),
				}

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "event created fail",
			fields: func() fields {
				u := &mocks.EventServiceInterface{}
				u.On("Create", eventservice.EventCreateRequest{
					Name:     "event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				}).Return(nil, errors.New("create error"))
				return fields{
					eventservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"event name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewEventHandler(tt.fields.log, tt.fields.eventservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestEventHandler_Get(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "event found",
			fields: func() fields {
				event := &eventservice.Event{}
				eventservice := &mocks.EventServiceInterface{}
				eventservice.On("Get", uint(1)).
					Return(event, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: eventservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "event not found",
			fields: func() fields {
				eventservice := &mocks.EventServiceInterface{}
				eventservice.On("Get", uint(1)).
					Return(nil, errors.New("event not found"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: eventservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "id is not int",
			fields: func() fields {
				return fields{
					log: logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "one",
					},
				}

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewEventHandler(
				tt.fields.log,
				tt.fields.eventservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestEventHandler_Update(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "event update success",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)

				u.On(
					"Update",
					event,
					eventservice.EventUpdateRequest{
						Name:     "new event name",
						Platform: []string{"platform 1", "platform 2"},
						Channel:  []string{"channel 1", "channel 2"},
					},
				).Return(&eventservice.Event{
					Name:     "new event name",
					Platform: []string{"platform 1", "platform 2"},
					Channel:  []string{"channel 1", "channel 2"},
				}, nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new event name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "event not found",
			fields: func() fields {
				u := &mocks.EventServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("event not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new event name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "body invalid",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new event name","platform":"platform 1","channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "event update fail",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)

				u.On(
					"Update",
					event,
					eventservice.EventUpdateRequest{
						Name:     "new event name",
						Platform: []string{"platform 1", "platform 2"},
						Channel:  []string{"channel 1", "channel 2"},
					},
				).Return(nil, errors.New("update fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new event name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "event id invalid",
			fields: func() fields {
				return fields{
					log: logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new event name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "one",
					},
				}

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewEventHandler(
				tt.fields.log,
				tt.fields.eventservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestEventHandler_Delete(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		eventservice eventservice.EventServiceInterface
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "event delete success",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)
				u.On("Delete", event).Return(nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				c.Set("event", &eventservice.Event{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "event not found",
			fields: func() fields {
				u := &mocks.EventServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("event not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				c.Set("event", &eventservice.Event{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "event delete fail",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)
				u.On("Delete", event).Return(errors.New("delete fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				c.Set("event", &eventservice.Event{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "event id invalid",
			fields: func() fields {
				return fields{
					log: logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "one",
					},
				}

				c.Set("event", &eventservice.Event{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "current event is incorrect",
			fields: func() fields {
				return fields{
					log: logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				c.Set("event", "1")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "cannot suicide",
			fields: func() fields {
				event := &eventservice.Event{}
				u := &mocks.EventServiceInterface{}

				u.On("Get", uint(1)).Return(event, nil)
				u.On("Delete", event).Return(nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					eventservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
				}

				c.Set("event", &eventservice.Event{Model: model.Model{ID: 1}})

				return args{c}
			}(),
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewEventHandler(
				tt.fields.log,
				tt.fields.eventservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
