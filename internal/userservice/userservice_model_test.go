package userservice

import (
	"testing"
)

func TestUser_SetPassword(t *testing.T) {
	type fields struct {
		Username    string
		password    string
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
			u := &User{
				Username:    tt.fields.Username,
				password:    tt.fields.password,
				DisplayName: tt.fields.DisplayName,
			}
			u.SetPassword(tt.args.password)
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	type fields struct {
		Username    string
		password    string
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
				password: "$2a$10$B2r2aAadfOjIFCyOg9HLS.TyE6RYWViuZj78p6zRvfJcIGjmWPA/m",
			},
			args: args{
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "password incorrect",
			fields: fields{
				password: "$2a$10$B2r2aAadfOjIFCyOg9HLS.TyE6RYWViuZj78p6zRvfJcIGjmWPA/m",
			},
			args: args{
				password: "Password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Username:    tt.fields.Username,
				password:    tt.fields.password,
				DisplayName: tt.fields.DisplayName,
			}
			if err := u.ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("User.ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
