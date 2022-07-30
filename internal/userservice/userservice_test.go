package userservice

import (
	"reflect"
	"testing"

	"github.com/oleiade/reflections"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want UserService
	}{
		{
			name: "new user",
			want: UserService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	type args struct {
		r UserCreateRequest
	}
	tests := []struct {
		name                string
		u                   UserService
		args                args
		want                map[string]interface{}
		wantErr             bool
		passwordValidateErr bool
	}{
		{
			name: "user created",
			u:    UserService{},
			args: args{
				r: UserCreateRequest{
					Username:    "admin",
					Password:    "password",
					DisplayName: "Administrator",
				},
			},
			want: map[string]interface{}{
				"Username":    "admin",
				"DisplayName": "Administrator",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserService{}
			got, err := u.Create(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if func() bool {
				for f, v := range tt.want {
					if k, err := reflections.GetField(tt.want, f); err != nil || k != v {
						return false
					}
				}

				return true
			}() {
				t.Errorf("UserService.Create() = %v, want %v", got, tt.want)
			}

			if (got.ValidatePassword(tt.args.r.Password) != nil) != tt.passwordValidateErr {
				t.Errorf("UserService.User.PasswordValidation() = %v, want %v", got.ValidatePassword(tt.args.r.Password), tt.passwordValidateErr)
			}
		})
	}
}
