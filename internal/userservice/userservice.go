package userservice

import "gorm.io/gorm"

type UserServiceDatabaseInterface interface {
	Create(value interface{}) (tx *gorm.DB)
}

type UserService struct {
	db UserServiceDatabaseInterface
}

type UserServiceInterface interface {
	Create(r UserCreateRequest) (*User, error)
}

func New(db UserServiceDatabaseInterface) UserService {
	return UserService{db}
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
