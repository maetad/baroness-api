package userservice

import (
	"github.com/maetad/baroness-api/internal/database"
	"github.com/maetad/baroness-api/internal/model"
)

type UserService struct {
	db database.DatabaseInterface
}

type UserServiceInterface interface {
	List() ([]model.UserInterface, error)
	Create(r UserCreateRequest) (model.UserInterface, error)
	Get(id uint) (model.UserInterface, error)
	GetByUsername(username string) (model.UserInterface, error)
	Update(user model.UserInterface, r UserUpdateRequest) (model.UserInterface, error)
	Delete(user model.UserInterface) error
}

func New(db database.DatabaseInterface) UserServiceInterface {
	return UserService{db}
}

func (s UserService) List() ([]model.UserInterface, error) {
	var users []model.User
	if result := s.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}

	var u = make([]model.UserInterface, len(users))

	for i := 0; i < len(users); i++ {
		u[i] = &users[i]
	}

	return u, nil
}

func (s UserService) Create(r UserCreateRequest) (model.UserInterface, error) {
	user := &model.User{
		Username:    r.Username,
		DisplayName: r.DisplayName,
	}

	user.SetPassword(r.Password)

	if result := s.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s UserService) Get(id uint) (model.UserInterface, error) {
	user := &model.User{}

	if result := s.db.First(user, id); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s UserService) GetByUsername(username string) (model.UserInterface, error) {
	user := &model.User{
		Username: username,
	}

	if result := s.db.First(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s UserService) Update(user model.UserInterface, r UserUpdateRequest) (model.UserInterface, error) {
	u := user.(*model.User)
	u.DisplayName = r.DisplayName
	if r.Password != "" {
		u.SetPassword(r.Password)
	}

	if result := s.db.Save(u); result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (s UserService) Delete(user model.UserInterface) error {
	return s.db.Delete(user).Error
}
