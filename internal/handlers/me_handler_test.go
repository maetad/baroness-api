package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pakkaparn/no-idea-api/internal/handlers"
	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/sirupsen/logrus"
)

func TestMeHandler_Get(t *testing.T) {
	type fields struct {
		log *logrus.Entry
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

				c.Params = gin.Params{
					{
						Key:   "id",
						Value: "1",
					},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewMeHandler(tt.fields.log)
			h.Get(tt.args.c)

			if tt.args.c.Writer.Status() != tt.want {
				t.Errorf("Get() = %v, want %v", tt.args.c.Writer.Status(), tt.want)
			}
		})
	}
}
