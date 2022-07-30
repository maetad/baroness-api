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
	"github.com/golang-jwt/jwt/v4"
	"github.com/maetad/baroness-api/internal/config"
	"github.com/maetad/baroness-api/internal/handlers"
	"github.com/maetad/baroness-api/internal/model"
	"github.com/maetad/baroness-api/internal/services/authservice"
	"github.com/maetad/baroness-api/internal/services/userservice"
	"github.com/maetad/baroness-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestNewAuthHandler(t *testing.T) {
	type args struct {
		log         *logrus.Entry
		options     config.Options
		authservice authservice.AuthServiceInterface
		userservice userservice.UserServiceInterface
	}
	tests := []struct {
		name string
		args args
		want *handlers.AuthHandler
	}{
		{
			name: "create auth handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlers.NewAuthHandler(tt.args.log, tt.args.options, tt.args.authservice, tt.args.userservice); reflect.TypeOf(got) != reflect.TypeOf(&handlers.AuthHandler{}) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log         *logrus.Entry
		authservice authservice.AuthServiceInterface
		userservice userservice.UserServiceInterface
		options     config.Options
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
			name: "invalid payload",
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"user":"username","pass":"password"}`)),
				}

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "user not found",
			fields: func() fields {
				userservice := &mocks.UserServiceInterface{}
				userservice.On("GetByUsername", mock.AnythingOfType("string")).
					Return(nil, errors.New("user not found"))

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password"}`)),
				}

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "user password incorrect",
			fields: func() fields {
				user := &mocks.UserInterface{}
				userservice := &mocks.UserServiceInterface{}
				user.On("ValidatePassword", mock.AnythingOfType("string")).
					Return(errors.New("password incorrect"))

				userservice.On("GetByUsername", mock.AnythingOfType("string")).
					Return(user, nil)

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password"}`)),
				}

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "generate token fail",
			fields: func() fields {
				user := &model.User{
					Username:    "admin",
					DisplayName: "administrator",
				}
				user.SetPassword("password")
				userservice := &mocks.UserServiceInterface{}
				authservice := &mocks.AuthServiceInterface{}

				userservice.On("GetByUsername", mock.AnythingOfType("string")).
					Return(user, nil)

				authservice.On("GenerateToken", mock.Anything, mock.Anything).
					Return("", errors.New("generate token fail"))

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					userservice: userservice,
					authservice: authservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password"}`)),
				}

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
		{
			name: "logged in",
			fields: func() fields {
				user := &model.User{
					Username:    "admin",
					DisplayName: "administrator",
				}
				user.SetPassword("password")

				userservice := &mocks.UserServiceInterface{}
				authservice := &mocks.AuthServiceInterface{}

				userservice.On("GetByUsername", mock.AnythingOfType("string")).
					Return(user, nil)

				authservice.On("GenerateToken", mock.Anything, mock.Anything).
					Return("token", nil)

				f := fields{
					userservice: userservice,
					authservice: authservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
					Body:   io.NopCloser(strings.NewReader(`{"username":"username","password":"password"}`)),
				}

				return args{c}
			}(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewAuthHandler(
				tt.fields.log,
				tt.fields.options,
				tt.fields.authservice,
				tt.fields.userservice,
			)
			h.Login(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Login() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestAuthHandler_Authorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		log         *logrus.Entry
		options     config.Options
		authservice authservice.AuthServiceInterface
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
			name: "token invalid",
			fields: func() fields {
				authservice := &mocks.AuthServiceInterface{}

				authservice.On("ParseToken", "jwttoken").
					Return(nil, errors.New("cannot parse"))

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					authservice: authservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Request.Header.Set("Authorization", "Bearer jwttoken")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "claim id not exists",
			fields: func() fields {
				authservice := &mocks.AuthServiceInterface{}

				authservice.On("ParseToken", "jwttoken").
					Return(jwt.MapClaims{}, nil)

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					authservice: authservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Request.Header.Set("Authorization", "Bearer jwttoken")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "claim username invalid",
			fields: func() fields {
				authservice := &mocks.AuthServiceInterface{}

				authservice.On("ParseToken", "jwttoken").
					Return(jwt.MapClaims{"username": nil}, nil)

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					authservice: authservice,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Request.Header.Set("Authorization", "Bearer jwttoken")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			fields: func() fields {
				authservice := &mocks.AuthServiceInterface{}

				authservice.On("ParseToken", "jwttoken").
					Return(jwt.MapClaims{"username": "admin"}, nil)

				u := &mocks.UserServiceInterface{}
				u.On("GetByUsername", "admin").
					Return(nil, errors.New("user not found"))

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					authservice: authservice,
					userservice: u,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Request.Header.Set("Authorization", "Bearer jwttoken")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "token valid",
			fields: func() fields {
				authservice := &mocks.AuthServiceInterface{}

				authservice.On("ParseToken", "jwttoken").
					Return(jwt.MapClaims{"username": "admin"}, nil)

				u := &mocks.UserServiceInterface{}
				u.On("GetByUsername", "admin").
					Return(&model.User{}, nil)

				f := fields{
					log:         logrus.WithContext(context.TODO()),
					authservice: authservice,
					userservice: u,
				}

				return f
			}(),
			args: func() args {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = &http.Request{
					URL:    &url.URL{},
					Header: make(http.Header),
				}

				c.Request.Header.Set("Authorization", "Bearer jwttoken")

				return args{c}
			}(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewAuthHandler(
				tt.fields.log,
				tt.fields.options,
				tt.fields.authservice,
				tt.fields.userservice,
			)
			h.Authorize(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Authorize() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
