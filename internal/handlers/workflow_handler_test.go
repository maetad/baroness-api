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
	"github.com/maetad/baroness-api/internal/services/workflowservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestNewWorkflowHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.WorkflowHandler
	}{
		{
			name: "create workflow handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewWorkflowHandler(tt.args.log, tt.args.workflowservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.WorkflowHandler{}) {
				t.Errorf("NewWorkflowHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
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
				workflows := make([]model.Workflow, 1)
				workflowservice := &mocks.WorkflowServiceInterface{}
				workflowservice.On("List", mock.Anything).
					Return(workflows, nil)
				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: workflowservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "listed fail",
			fields: func() fields {
				workflowservice := &mocks.WorkflowServiceInterface{}
				workflowservice.On("List", mock.Anything).
					Return(nil, errors.New("list fail"))
				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: workflowservice,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewWorkflowHandler(tt.fields.log, tt.fields.workflowservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestWorkflowHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
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
			name: "workflow created success",
			fields: func() fields {
				u := &mocks.WorkflowServiceInterface{}
				u.On("Create", workflowservice.WorkflowCreateRequest{
					Name: "workflow name",
				}, mock.Anything).Return(&model.Workflow{}, nil)
				return fields{
					workflowservice: u,
					log:             logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"workflow name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "workflow created fail invalid payload",
			fields: func() fields {
				u := &mocks.WorkflowServiceInterface{}
				u.On("Create", workflowservice.WorkflowCreateRequest{
					Name: "workflow name",
				}, mock.Anything).Return(&model.Workflow{}, nil)
				return fields{
					workflowservice: u,
					log:             logrus.WithContext(context.TODO()),
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
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "workflow created fail",
			fields: func() fields {
				u := &mocks.WorkflowServiceInterface{}
				u.On("Create", workflowservice.WorkflowCreateRequest{
					Name: "workflow name",
				}, mock.Anything).Return(nil, errors.New("create error"))
				return fields{
					workflowservice: u,
					log:             logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"workflow name","platform":["platform 1", "platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewWorkflowHandler(tt.fields.log, tt.fields.workflowservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestWorkflowHandler_Get(t *testing.T) {
	type fields struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
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
			name: "workflow found",
			fields: func() fields {
				workflow := &model.Workflow{}
				workflowservice := &mocks.WorkflowServiceInterface{}
				workflowservice.On("Get", uint(1)).
					Return(workflow, nil)
				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: workflowservice,
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
						Key:   "workflow_id",
						Value: "1",
					},
				}

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "workflow not found",
			fields: func() fields {
				workflowservice := &mocks.WorkflowServiceInterface{}
				workflowservice.On("Get", uint(1)).
					Return(nil, errors.New("workflow not found"))
				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: workflowservice,
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
						Key:   "workflow_id",
						Value: "1",
					},
				}

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
						Key:   "workflow_id",
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
			h := handlers.NewWorkflowHandler(
				tt.fields.log,
				tt.fields.workflowservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestWorkflowHandler_Update(t *testing.T) {
	type fields struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
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
			name: "workflow update success",
			fields: func() fields {
				workflow := &model.Workflow{}
				u := &mocks.WorkflowServiceInterface{}

				u.On("Get", uint(1)).Return(workflow, nil)

				u.On(
					"Update",
					workflow,
					workflowservice.WorkflowUpdateRequest{
						Name: "new workflow name",
					},
					mock.Anything,
				).Return(&model.Workflow{
					Name: "new workflow name",
				}, nil)

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new workflow name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
		{
			name: "workflow not found",
			fields: func() fields {
				u := &mocks.WorkflowServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("workflow not found"))

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new workflow name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

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
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "workflow update fail",
			fields: func() fields {
				workflow := &model.Workflow{}
				u := &mocks.WorkflowServiceInterface{}

				u.On("Get", uint(1)).Return(workflow, nil)

				u.On(
					"Update",
					workflow,
					workflowservice.WorkflowUpdateRequest{
						Name: "new workflow name",
					},
					mock.Anything,
				).Return(nil, errors.New("update fail"))

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new workflow name","platform":["platform 1","platform 2"],"channel":["channel 1", "channel 2"]}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "workflow id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new workflow name"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "workflow_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})
				c.Set("event", &model.Event{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewWorkflowHandler(
				tt.fields.log,
				tt.fields.workflowservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestWorkflowHandler_Delete(t *testing.T) {
	type fields struct {
		log             *logrus.Entry
		workflowservice workflowservice.WorkflowServiceInterface
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
			name: "workflow delete success",
			fields: func() fields {
				workflow := &model.Workflow{}
				u := &mocks.WorkflowServiceInterface{}

				u.On("Get", uint(1)).Return(workflow, nil)
				u.On("Delete", workflow, mock.Anything).Return(nil)

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
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
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "workflow not found",
			fields: func() fields {
				u := &mocks.WorkflowServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("workflow not found"))

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
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
						Key:   "workflow_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "workflow id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "workflow_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "workflow delete fail",
			fields: func() fields {
				workflow := &model.Workflow{}
				u := &mocks.WorkflowServiceInterface{}

				u.On("Get", uint(1)).Return(workflow, nil)
				u.On("Delete", workflow, mock.Anything).Return(errors.New("delete fail"))

				return fields{
					log:             logrus.WithContext(context.TODO()),
					workflowservice: u,
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
						Key:   "workflow_id",
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
			h := handlers.NewWorkflowHandler(
				tt.fields.log,
				tt.fields.workflowservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
