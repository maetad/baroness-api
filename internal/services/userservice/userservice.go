package userservice

import "gorm.io/gorm"

type UserServiceDatabaseInterface interface {
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
}

type UserService struct {
	db UserServiceDatabaseInterface
}

type UserServiceInterface interface {
	List() ([]UserInterface, error)
	Create(r UserCreateRequest) (UserInterface, error)
	Get(id uint) (UserInterface, error)
	GetByUsername(username string) (UserInterface, error)
}

func New(db UserServiceDatabaseInterface) UserServiceInterface {
	return UserService{db}
}

func (s UserService) List() ([]UserInterface, error) {
	var users []User
	if result := s.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}

	var u = make([]UserInterface, len(users))

	for i := 0; i < len(users); i++ {
		u[i] = &users[i]
	}

	return u, nil
}

func (s UserService) Create(r UserCreateRequest) (UserInterface, error) {
	user := &User{
		Username:    r.Username,
		DisplayName: r.DisplayName,
	}

	user.SetPassword(r.Password)

	if result := s.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s UserService) Get(id uint) (UserInterface, error) {
	user := &User{}

	if result := s.db.First(user, id); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s UserService) GetByUsername(username string) (UserInterface, error) {
	user := &User{
		Username: username,
	}

	if result := s.db.First(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
