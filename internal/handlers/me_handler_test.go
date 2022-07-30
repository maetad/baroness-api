package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/handlers"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/pakkaparn/no-idea-api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestMeHandler_Get(t *testing.T) {
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
			name: "get me success",
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusOK,
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

				c.Set("user", "1")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewMeHandler(tt.fields.log, tt.fields.userservice)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}

func TestMeHandler_Update(t *testing.T) {
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
			name: "me update success",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}

				u.On(
					"Update",
					mock.AnythingOfType("*userservice.User"),
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusOK,
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
					Body:   io.NopCloser(strings.NewReader(`{"display_name":"display_name","password":"password"}`)),
				}

				c.Set("user", "1")

				return args{c}
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "body invalid",
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
					Body:   io.NopCloser(strings.NewReader(`{"display___name":"display_name","password":"password"}`)),
				}

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "me update fail",
			fields: func() fields {
				u := &mocks.UserServiceInterface{}

				u.On(
					"Update",
					mock.AnythingOfType("*userservice.User"),
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

				c.Set("user", &userservice.User{})

				return args{c}
			}(),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewMeHandler(tt.fields.log, tt.fields.userservice)
			h.Update(tt.args.c)
		})

		if tt.args.c.Writer.Status() != tt.want {
			t.Errorf("Update() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
		}
	}
}
