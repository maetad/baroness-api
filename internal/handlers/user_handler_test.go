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
	"github.com/pakkaparn/no-idea-api/internal/handlers"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/pakkaparn/no-idea-api/mocks"
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
