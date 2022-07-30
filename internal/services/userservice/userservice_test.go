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
				user := got.(*userservice.User)

				if err := got.ValidatePassword(tt.args.r.Password); err != nil {
					t.Errorf("UserService.Create() password hashed invalid %v", err)
				}

				// ignore password difference
				user.Password = ""
				clone := tt.want
				clone.Password = ""

				if !reflect.DeepEqual(user, clone) {
					t.Errorf("UserService.Create() = %v, want %v", user, tt.want)
				}
			} else if got != nil {
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

			if tt.want != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UserService.GetByUsername() = %v, want %v", got, tt.want)
				}
			} else if got != nil {
				t.Errorf("UserService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_List(t *testing.T) {
	type fields struct {
		db userservice.UserServiceDatabaseInterface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []userservice.UserInterface
		wantErr bool
	}{
		{
			name: "listed success",
			fields: func() fields {
				db := &mocks.UserServiceDatabaseInterface{}
				db.On("Find", mock.Anything).
					Return(&gorm.DB{
						Error: nil,
					})

				return fields{db}
			}(),
			want: []userservice.UserInterface{},
		},
		{
			name: "listed fail",
			fields: func() fields {
				db := &mocks.UserServiceDatabaseInterface{}
				db.On("Find", mock.Anything).
					Return(&gorm.DB{
						Error: errors.New("find error"),
					})

				return fields{db}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := userservice.New(tt.fields.db)
			got, err := s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
