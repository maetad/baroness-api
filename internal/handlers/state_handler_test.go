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
	"github.com/maetad/baroness-api/internal/services/stateservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestNewStateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.StateHandler
	}{
		{
			name: "create state handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewStateHandler(tt.args.log, tt.args.stateservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.StateHandler{}) {
				t.Errorf("NewStateHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
				states := make([]model.State, 1)
				stateservice := &mocks.StateServiceInterface{}
				stateservice.On("List", mock.Anything).
					Return(states, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: stateservice,
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
				stateservice := &mocks.StateServiceInterface{}
				stateservice.On("List", mock.Anything).
					Return(nil, errors.New("list fail"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: stateservice,
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
			h := handlers.NewStateHandler(tt.fields.log, tt.fields.stateservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestStateHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
			name: "state created success",
			fields: func() fields {
				u := &mocks.StateServiceInterface{}
				u.On("Create", stateservice.StateCreateRequest{
					Name: "state name",
				}, mock.Anything).Return(&model.State{}, nil)
				return fields{
					stateservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"state name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "state created fail invalid payload",
			fields: func() fields {
				u := &mocks.StateServiceInterface{}
				u.On("Create", stateservice.StateCreateRequest{
					Name: "state name",
				}, mock.Anything).Return(&model.State{}, nil)
				return fields{
					stateservice: u,
					log:          logrus.WithContext(context.TODO()),
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
			name: "state created fail",
			fields: func() fields {
				u := &mocks.StateServiceInterface{}
				u.On("Create", stateservice.StateCreateRequest{
					Name: "state name",
				}, mock.Anything).Return(nil, errors.New("create error"))
				return fields{
					stateservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"state name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
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
			h := handlers.NewStateHandler(tt.fields.log, tt.fields.stateservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestStateHandler_Get(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
			name: "state found",
			fields: func() fields {
				state := &model.State{}
				stateservice := &mocks.StateServiceInterface{}
				stateservice.On("Get", uint(1)).
					Return(state, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: stateservice,
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
						Key:   "state_id",
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
			name: "state not found",
			fields: func() fields {
				stateservice := &mocks.StateServiceInterface{}
				stateservice.On("Get", uint(1)).
					Return(nil, errors.New("state not found"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: stateservice,
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
						Key:   "state_id",
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
						Key:   "state_id",
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
			h := handlers.NewStateHandler(
				tt.fields.log,
				tt.fields.stateservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestStateHandler_Update(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
			name: "state update success",
			fields: func() fields {
				state := &model.State{}
				u := &mocks.StateServiceInterface{}

				u.On("Get", uint(1)).Return(state, nil)

				u.On(
					"Update",
					state,
					stateservice.StateUpdateRequest{
						Name: "new state name",
					},
					mock.Anything,
				).Return(&model.State{
					Name: "new state name",
				}, nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new state name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "state_id",
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
			name: "state not found",
			fields: func() fields {
				u := &mocks.StateServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("state not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new state name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "state_id",
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
						Key:   "state_id",
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
			name: "state update fail",
			fields: func() fields {
				state := &model.State{}
				u := &mocks.StateServiceInterface{}

				u.On("Get", uint(1)).Return(state, nil)

				u.On(
					"Update",
					state,
					stateservice.StateUpdateRequest{
						Name: "new state name",
					},
					mock.Anything,
				).Return(nil, errors.New("update fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new state name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "state_id",
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
			name: "state id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new state name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "state_id",
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
			h := handlers.NewStateHandler(
				tt.fields.log,
				tt.fields.stateservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestStateHandler_Delete(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
			name: "state delete success",
			fields: func() fields {
				state := &model.State{}
				u := &mocks.StateServiceInterface{}

				u.On("Get", uint(1)).Return(state, nil)
				u.On("Delete", state, mock.Anything).Return(nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
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
						Key:   "state_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "state not found",
			fields: func() fields {
				u := &mocks.StateServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("state not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
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
						Key:   "state_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "state id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "state_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "state delete fail",
			fields: func() fields {
				state := &model.State{}
				u := &mocks.StateServiceInterface{}

				u.On("Get", uint(1)).Return(state, nil)
				u.On("Delete", state, mock.Anything).Return(errors.New("delete fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					stateservice: u,
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
						Key:   "state_id",
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
			h := handlers.NewStateHandler(
				tt.fields.log,
				tt.fields.stateservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestStateHandler_Show(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		stateservice stateservice.StateServiceInterface
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
			name: "state show success",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("state", &model.State{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewStateHandler(tt.fields.log, tt.fields.stateservice)
			h.Show(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Show() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
