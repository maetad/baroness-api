package userservice

type UserService struct{}

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

	return user, nil
}
