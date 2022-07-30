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
	"github.com/maetad/baroness-api/internal/services/userservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
)

func TestNewUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.UserHandler
	}{
		{
			name: "create user handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewUserHandler(tt.args.log, tt.args.userservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.UserHandler{}) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
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
				users := make([]userservice.UserInterface, 1)
				userservice := &mocks.UserServiceInterface{}
				userservice.On("List").
					Return(users, nil)
				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
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
				userservice := &mocks.UserServiceInterface{}
				userservice.On("List").
					Return(nil, errors.New("list fail"))
				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
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
			h := handlers.NewUserHandler(tt.fields.log, tt.fields.userservice)
			h.List(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("List() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
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
			name: "user created success",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}
				u.On("Create", userservice.UserCreateRequest{
					Username:    "username",
					Password:    "password",
					DisplayName: "Adminstrator",
				}).Return(&userservice.User{}, nil)
				return fields{
					userservice: u,
					log:         logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password","display_name":"Adminstrator"}`)),
				}

				return args{c}
			}(),
			want: http.StatusCreated,
		},
		{
			name: "user created fail invalid payload",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}
				u.On("Create", userservice.UserCreateRequest{
					Username:    "username",
					Password:    "password",
					DisplayName: "Adminstrator",
				}).Return(&userservice.User{}, nil)
				return fields{
					userservice: u,
					log:         logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"","display_name":"Adminstrator"}`)),
				}

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "user created fail",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}
				u.On("Create", userservice.UserCreateRequest{
					Username:    "username",
					Password:    "password",
					DisplayName: "Adminstrator",
				}).Return(nil, errors.New("create error"))
				return fields{
					userservice: u,
					log:         logrus.WithContext(context.TODO()),
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password","display_name":"Adminstrator"}`)),
				}

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewUserHandler(tt.fields.log, tt.fields.userservice)
			h.Create(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Create() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
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
			name: "user found",
			fields: func() fields {
				user := &userservice.User{}
				userservice := &mocks.UserServiceInterface{}
				userservice.On("Get", uint(1)).
					Return(user, nil)
				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
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
			name: "user not found",
			fields: func() fields {
				userservice := &mocks.UserServiceInterface{}
				userservice.On("Get", uint(1)).
					Return(nil, errors.New("user not found"))
				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
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
			h := handlers.NewUserHandler(
				tt.fields.log,
				tt.fields.userservice,
			)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
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
			name: "user update success",
			fields: func() fields {
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)

				u.On(
					"Update",
					user,
					userservice.UserUpdateRequest{DisplayName: "display_name", Password: "password"},
				).Return(&userservice.User{
					DisplayName: "display_name",
					Password:    "password",
				}, nil)

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"display_name":"display_name","password":"password"}`)),
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
			name: "user not found",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("user not found"))

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"display_name":"display_name","password":"password"}`)),
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
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"display___name":"display_name","password":"password"}`)),
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
			name: "user update fail",
			fields: func() fields {
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)

				u.On(
					"Update",
					user,
					userservice.UserUpdateRequest{DisplayName: "display_name", Password: "password"},
				).Return(nil, errors.New("update fail"))

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
				}
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"display_name":"display_name","password":"password"}`)),
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
			name: "user id invalid",
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
					Body:   io.NopCloser(strings.NewReader(`{"display_name":"display_name","password":"password"}`)),
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
			h := handlers.NewUserHandler(
				tt.fields.log,
				tt.fields.userservice,
			)
			h.Update(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		userservice userservice.UserServiceInterface
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
			name: "user delete success",
			fields: func() fields {
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)
				u.On("Delete", user).Return(nil)

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusNoContent,
		},
		{
			name: "user not found",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}
				u.On("Get", uint(1)).Return(nil, errors.New("user not found"))

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "user delete fail",
			fields: func() fields {
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)
				u.On("Delete", user).Return(errors.New("delete fail"))

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "user id invalid",
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "current user is incorrect",
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

				c.Set("user", "1")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "cannot suicide",
			fields: func() fields {
				user := &userservice.User{}
				u := &mocks.UserServiceInterface{}

				u.On("Get", uint(1)).Return(user, nil)
				u.On("Delete", user).Return(nil)

				return fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: u,
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

				c.Set("user", &userservice.User{Model: model.Model{ID: 1}})

				return args{c}
			}(),
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewUserHandler(
				tt.fields.log,
				tt.fields.userservice,
			)
			h.Delete(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Delete() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
