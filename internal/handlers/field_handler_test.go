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
	"github.com/maetad/baroness-api/internal/services/fieldservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestNewFieldHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.FieldHandler
	}{
		{
			name: "create field handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewFieldHandler(tt.args.log, tt.args.fieldservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.FieldHandler{}) {
				t.Errorf("NewFieldHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
				f := make([]model.Field, 1)
				fieldservice := &mocks.FieldServiceInterface{}
				fieldservice.On("List", mock.Anything).
					Return(f, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: fieldservice,
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
				fieldservice := &mocks.FieldServiceInterface{}
				fieldservice.On("List", mock.Anything).
					Return(nil, errors.New("list fail"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: fieldservice,
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
			h := handlers.NewFieldHandler(tt.fields.log, tt.fields.fieldservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestFieldHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
			name: "field created success",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Create", fieldservice.FieldCreateRequest{
					Name: "field name",
					Type: "text",
				}, mock.Anything).Return(&model.Field{}, nil)
				return fields{
					fieldservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"field name","type":"text"}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "field created fail invalid payload",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Create", fieldservice.FieldCreateRequest{
					Name: "field name",
				}, mock.Anything).Return(&model.Field{}, nil)
				return fields{
					fieldservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name": 1,"type":"text"}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "field created fail type is not in enum",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Create", fieldservice.FieldCreateRequest{
					Name: "field name",
				}, mock.Anything).Return(&model.Field{}, nil)
				return fields{
					fieldservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name": "field name","type":"string"}`)),
				}

				c.Set("user", &model.User{})
				c.Set("workflow", &model.Workflow{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "field created fail",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Create", fieldservice.FieldCreateRequest{
					Name: "field name",
					Type: "text",
				}, mock.Anything).Return(nil, errors.New("create error"))
				return fields{
					fieldservice: u,
					log:          logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"field name","type":"text"}`)),
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
			h := handlers.NewFieldHandler(tt.fields.log, tt.fields.fieldservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestFieldHandler_Get(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
			name: "field found",
			fields: func() fields {
				field := &model.Field{}
				fieldservice := &mocks.FieldServiceInterface{}
				fieldservice.On("Get", uint(1)).
					Return(field, nil)
				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: fieldservice,
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
						Key:   "field_id",
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
			name: "field not found",
			fields: func() fields {
				fieldservice := &mocks.FieldServiceInterface{}
				fieldservice.On("Get", uint(1)).
					Return(nil, errors.New("field not found"))
				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: fieldservice,
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
						Key:   "field_id",
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
						Key:   "field_id",
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
			h := handlers.NewFieldHandler(
				tt.fields.log,
				tt.fields.fieldservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestFieldHandler_Update(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
			name: "field update success",
			fields: func() fields {
				field := &model.Field{}
				u := &mocks.FieldServiceInterface{}

				u.On("Get", uint(1)).Return(field, nil)

				u.On(
					"Update",
					field,
					fieldservice.FieldUpdateRequest{
						Name: "new field name",
						Type: "text",
					},
					mock.Anything,
				).Return(&model.Field{
					Name: "new field name",
				}, nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new field name","type":"text"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
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
			name: "field not found",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("field not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new field name","type":"text"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
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
			name: "field updated fail invalid payload",
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
						Key:   "field_id",
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
			name: "field updated fail type is not enum",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name": "new field name","type": "string"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
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
			name: "field update fail",
			fields: func() fields {
				field := &model.Field{}
				u := &mocks.FieldServiceInterface{}

				u.On("Get", uint(1)).Return(field, nil)

				u.On(
					"Update",
					field,
					fieldservice.FieldUpdateRequest{
						Name: "new field name",
						Type: "text",
					},
					mock.Anything,
				).Return(nil, errors.New("update fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new field name","type":"text"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
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
			name: "field id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"name":"new field name","type":"text"}`)),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
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
			h := handlers.NewFieldHandler(
				tt.fields.log,
				tt.fields.fieldservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestFieldHandler_Delete(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
			name: "field delete success",
			fields: func() fields {
				field := &model.Field{}
				u := &mocks.FieldServiceInterface{}

				u.On("Get", uint(1)).Return(field, nil)
				u.On("Delete", field, mock.Anything).Return(nil)

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
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
						Key:   "field_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "field not found",
			fields: func() fields {
				u := &mocks.FieldServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("field not found"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
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
						Key:   "field_id",
						Value: "1",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "field id invalid",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Params = gin.Params{
					{
						Key:   "field_id",
						Value: "one",
					},
				}

				c.Set("user", &model.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "field delete fail",
			fields: func() fields {
				field := &model.Field{}
				u := &mocks.FieldServiceInterface{}

				u.On("Get", uint(1)).Return(field, nil)
				u.On("Delete", field, mock.Anything).Return(errors.New("delete fail"))

				return fields{
					log:          logrus.WithContext(context.TODO()),
					fieldservice: u,
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
						Key:   "field_id",
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
			h := handlers.NewFieldHandler(
				tt.fields.log,
				tt.fields.fieldservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestFieldHandler_Show(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		fieldservice fieldservice.FieldServiceInterface
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
			name: "field show success",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Set("field", &model.Field{})

				return args{c}
			}(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewFieldHandler(tt.fields.log, tt.fields.fieldservice)
			h.Show(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Show() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
