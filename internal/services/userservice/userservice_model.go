package userservice

import (
	"github.com/pakkaparn/no-idea-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	SetPassword(password string)
	ValidatePassword(password string) error
}

type User struct {
	model.Model
	Username    string `json:"username"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
}

func (u *User) SetPassword(password string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u.Password = string(hashed)
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) GetClaims() map[string]interface{} {
	return map[string]interface{}{
		"username":     u.Username,
		"display_name": u.DisplayName,
	}
}
