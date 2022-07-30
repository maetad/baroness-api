package userservice

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username    string
	password    string
	DisplayName string
}

func (u *User) SetPassword(password string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u.password = string(hashed)
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
}
