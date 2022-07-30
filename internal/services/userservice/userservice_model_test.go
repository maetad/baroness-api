package userservice_test

import (
	"testing"

	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
)

func TestUser_SetPassword(t *testing.T) {
	type fields struct {
		Username    string
		DisplayName string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "set password",
			fields: fields{},
			args: args{
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userservice.User{
				Username:    tt.fields.Username,
				DisplayName: tt.fields.DisplayName,
			}
			u.SetPassword(tt.args.password)
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	type fields struct {
		Username    string
		Password    string
		DisplayName string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "password correct",
			fields: fields{
				Password: "$2a$10$B2r2aAadfOjIFCyOg9HLS.TyE6RYWViuZj78p6zRvfJcIGjmWPA/m",
			},
			args: args{
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "password incorrect",
			fields: fields{
				Password: "$2a$10$B2r2aAadfOjIFCyOg9HLS.TyE6RYWViuZj78p6zRvfJcIGjmWPA/m",
			},
			args: args{
				password: "Password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userservice.User{
				Username:    tt.fields.Username,
				Password:    tt.fields.Password,
				DisplayName: tt.fields.DisplayName,
			}
			if err := u.ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("User.ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
