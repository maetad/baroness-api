package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/handlers"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/pakkaparn/no-idea-api/mocks"
	"github.com/sirupsen/logrus"
)

func TestNewUserHandler(t *testing.T) {
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
					log:         &logrus.Entry{},
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
					log:         &logrus.Entry{},
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
		})
	}
}
