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
	"github.com/maetad/baroness-api/internal/services/transitionservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestNewTransitionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.TransitionHandler
	}{
		{
			name: "create transition handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewTransitionHandler(tt.args.log, tt.args.transitionservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.TransitionHandler{}) {
				t.Errorf("NewTransitionHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
				transitions := make([]model.Transition, 1)
				transitionservice := &mocks.TransitionServiceInterface{}
				transitionservice.On("List", mock.Anything).
					Return(transitions, nil)
				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: transitionservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "listed fail",
			fields: func() fields {
				transitionservice := &mocks.TransitionServiceInterface{}
				transitionservice.On("List", mock.Anything).
					Return(nil, errors.New("list fail"))
				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: transitionservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(tt.fields.log, tt.fields.transitionservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestTransitionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
			name: "transition created success",
			fields: func() fields {
				u := &mocks.TransitionServiceInterface{}
				u.On("Create", transitionservice.TransitionCreateRequest{
					Name: "transition name",
				}, mock.Anything).Return(&model.Transition{}, nil)
				return fields{
					transitionservice: u,
					log:               logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"transition name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "transition created fail invalid payload",
			fields: func() fields {
				u := &mocks.TransitionServiceInterface{}
				u.On("Create", transitionservice.TransitionCreateRequest{
					Name: "transition name",
				}, mock.Anything).Return(&model.Transition{}, nil)
				return fields{
					transitionservice: u,
					log:               logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name": 1`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "transition created fail",
			fields: func() fields {
				u := &mocks.TransitionServiceInterface{}
				u.On("Create", transitionservice.TransitionCreateRequest{
					Name: "transition name",
				}, mock.Anything).Return(nil, errors.New("create error"))
				return fields{
					transitionservice: u,
					log:               logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"transition name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(tt.fields.log, tt.fields.transitionservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestTransitionHandler_Get(t *testing.T) {
	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
			name: "transition found",
			fields: func() fields {
				transition := &model.Transition{}
				transitionservice := &mocks.TransitionServiceInterface{}
				transitionservice.On("Get", uint(1)).
					Return(transition, nil)
				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: transitionservice,
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
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "transition not found",
			fields: func() fields {
				transitionservice := &mocks.TransitionServiceInterface{}
				transitionservice.On("Get", uint(1)).
					Return(nil, errors.New("transition not found"))
				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: transitionservice,
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
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "id is not int",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(
				tt.fields.log,
				tt.fields.transitionservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestTransitionHandler_Update(t *testing.T) {
	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
			name: "transition update success",
			fields: func() fields {
				transition := &model.Transition{}
				u := &mocks.TransitionServiceInterface{}

				u.On("Get", uint(1)).Return(transition, nil)

				u.On(
					"Update",
					transition,
					transitionservice.TransitionUpdateRequest{
						Name: "new transition name",
					},
					mock.Anything,
				).Return(&model.Transition{
					Name: "new transition name",
				}, nil)

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new transition name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "transition not found",
			fields: func() fields {
				u := &mocks.TransitionServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("transition not found"))

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new transition name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "body invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name": 1}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "transition update fail",
			fields: func() fields {
				transition := &model.Transition{}
				u := &mocks.TransitionServiceInterface{}

				u.On("Get", uint(1)).Return(transition, nil)

				u.On(
					"Update",
					transition,
					transitionservice.TransitionUpdateRequest{
						Name: "new transition name",
					},
					mock.Anything,
				).Return(nil, errors.New("update fail"))

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new transition name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "transition id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new transition name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(
				tt.fields.log,
				tt.fields.transitionservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestTransitionHandler_Delete(t *testing.T) {
	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
			name: "transition delete success",
			fields: func() fields {
				transition := &model.Transition{}
				u := &mocks.TransitionServiceInterface{}

				u.On("Get", uint(1)).Return(transition, nil)
				u.On("Delete", transition, mock.Anything).Return(nil)

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
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
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "transition not found",
			fields: func() fields {
				u := &mocks.TransitionServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("transition not found"))

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
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
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "transition id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "transition_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "transition delete fail",
			fields: func() fields {
				transition := &model.Transition{}
				u := &mocks.TransitionServiceInterface{}

				u.On("Get", uint(1)).Return(transition, nil)
				u.On("Delete", transition, mock.Anything).Return(errors.New("delete fail"))

				return fields{
					log:               logrus.WithContext(context.TODO()),
					transitionservice: u,
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
						Key:   "transition_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(
				tt.fields.log,
				tt.fields.transitionservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestTransitionHandler_Show(t *testing.T) {
	type fields struct {
		log               *logrus.Entry
		transitionservice transitionservice.TransitionServiceInterface
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
			name: "transition show success",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("transition", &model.Transition{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewTransitionHandler(tt.fields.log, tt.fields.transitionservice)
			h.Show(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Show() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
