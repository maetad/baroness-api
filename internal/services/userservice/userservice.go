package userservice

import "gorm.io/gorm"

type UserServiceDatabaseInterface interface {
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
}

type UserService struct {
	db UserServiceDatabaseInterface
}

type UserServiceInterface interface {
	Create(r UserCreateRequest) (UserInterface, error)
	GetByUsername(username string) (UserInterface, error)
}

func New(db UserServiceDatabaseInterface) UserServiceInterface {
	return UserService{db}
}

func (u UserService) Create(r UserCreateRequest) (UserInterface, error) {
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

func (u UserService) GetByUsername(username string) (UserInterface, error) {
	user := &User{
		Username: username,
	}

	if result := u.db.First(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
