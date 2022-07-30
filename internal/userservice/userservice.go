package userservice

import "gorm.io/gorm"

type UserService struct {
	db *gorm.DB
}

type UserServiceInterface interface {
	Create(r UserCreateRequest) (*User, error)
}

func New() UserService {
	return UserService{}
}

func (u UserService) Create(r UserCreateRequest) (*User, error) {
	user := &User{
		Username:    r.Username,
		DisplayName: r.DisplayName,
	}

	user.SetPassword(r.Password)

	if result := u.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
