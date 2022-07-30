package userservice

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInterface interface {
	SetPassword(password string)
	ValidatePassword(password string) error
}

type User struct {
	gorm.Model
	Username    string
	Password    string
	DisplayName string
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
