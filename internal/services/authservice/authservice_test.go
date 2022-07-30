package authservice_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pakkaparn/no-idea-api/internal/services/authservice"
	"github.com/pakkaparn/no-idea-api/mocks"
)

var claimer = &mocks.Claimer{}
var jwtPattern = `^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)`
var jwtRegex = regexp.MustCompile(jwtPattern)

func TestAuthService_GenerateToken(t *testing.T) {
	type fields struct {
		signingMethod *jwt.SigningMethodHMAC
		signingKey    interface{}
	}
	type args struct {
		c authservice.Claimer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "token generated",
			fields: fields{
				signingMethod: jwt.SigningMethodHS256,
				signingKey:    []byte("signing-key"),
			},
			args: args{claimer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claimer.Mock.ExpectedCalls = nil
			claimer.On("GetClaims").Return(map[string]interface{}{
				"username": "admin",
			})

			s := authservice.New(tt.fields.signingMethod, tt.fields.signingKey)
			got, err := s.GenerateToken(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !jwtRegex.Match([]byte(got)) {
				t.Errorf("AuthService.GenerateToken() = %v, want %v", got, jwtPattern)
			}

			claimer.AssertNumberOfCalls(t, "GetClaims", 1)
		})
	}
}

