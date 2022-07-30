package model_test

import (
	"reflect"
	"testing"

	"github.com/maetad/baroness-api/internal/model"
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
			u := &model.User{
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
			u := &model.User{
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

func TestUser_GetClaims(t *testing.T) {
	type fields struct {
		Model       model.Model
		Username    string
		Password    string
		DisplayName string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "user claims",
			fields: fields{
				Username:    "admin",
				DisplayName: "Administrator",
			},
			want: map[string]interface{}{
				"username":     "admin",
				"display_name": "Administrator",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &model.User{
				Model:       tt.fields.Model,
				Username:    tt.fields.Username,
				Password:    tt.fields.Password,
				DisplayName: tt.fields.DisplayName,
			}
			if got := u.GetClaims(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.GetClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}
