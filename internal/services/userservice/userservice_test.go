package userservice_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pakkaparn/no-idea-api/internal/services/userservice"
	"github.com/pakkaparn/no-idea-api/mocks"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var db = &mocks.UserServiceDatabaseInterface{}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "new user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userservice.New(db); reflect.ValueOf(got).Kind() != reflect.ValueOf(userservice.UserService{}).Kind() {
				t.Errorf("New() = %v, want %v", reflect.ValueOf(got).Kind(), reflect.ValueOf(userservice.UserService{}).Kind())
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	type fields struct {
		db userservice.UserServiceDatabaseInterface
	}
	type args struct {
		r userservice.UserCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *userservice.User
		wantErr bool
	}{
		{
			name: "user created",
			fields: fields{
				db: db,
			},
			args: args{
				r: userservice.UserCreateRequest{
					Username:    "admin",
					Password:    "password",
					DisplayName: "Administrator",
				},
			},
			want: &userservice.User{
				Username:    "admin",
				Password:    "$2a$10$EIbuP5hbywq0xp183mHeBe0cN6TO00FNK7sAZJGKXWr9V6A2pVLkS",
				DisplayName: "Administrator",
			},
		},
		{
			name: "user create fail",
			fields: fields{
				db: db,
			},
			args: args{
				r: userservice.UserCreateRequest{
					Username:    "admin",
					Password:    "password",
					DisplayName: "Administrator",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Mock.ExpectedCalls = nil
			db.On("Create", mock.AnythingOfType("*userservice.User")).
				Return(&gorm.DB{
					Error: func() error {
						if tt.wantErr {
							return errors.New("error")
						}

						return nil
					}(),
				})

			u := userservice.New(tt.fields.db)
			got, err := u.Create(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if got.Username != tt.want.Username {
					t.Errorf("UserService.Create() = %v, want %v", got.Username, tt.want.Username)
				}
				if got.DisplayName != tt.want.DisplayName {
					t.Errorf("UserService.Create() = %v, want %v", got.DisplayName, tt.want.DisplayName)
				}

				if err := got.ValidatePassword(tt.args.r.Password); err != nil {
					t.Errorf("UserService.Create() password hashed invalid %v", err)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	type fields struct {
		db userservice.UserServiceDatabaseInterface
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *userservice.User
		wantErr bool
	}{
		{
			name: "found",
			fields: fields{
				db: db,
			},
			args: args{
				username: "admin",
			},
			want: &userservice.User{
				Username: "admin",
			},
		},
		{
			name: "not found",
			fields: fields{
				db: db,
			},
			args: args{
				username: "admin",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Mock.ExpectedCalls = nil
			db.On("First", mock.AnythingOfType("*userservice.User")).
				Return(&gorm.DB{
					Error: func() error {
						if tt.wantErr {
							return errors.New("error")
						}

						return nil
					}(),
				})

			u := userservice.New(tt.fields.db)
			got, err := u.GetByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
